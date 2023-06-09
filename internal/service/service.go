package service

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/k4zb3k/project/internal/apperror"
	"github.com/k4zb3k/project/internal/models"
	"github.com/k4zb3k/project/internal/repository"
	"github.com/k4zb3k/project/pkg/logger"
	"github.com/twinj/uuid"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	Repository *repository.Repository
	Redis      *redis.Client
}

func NewService(repository *repository.Repository, redis *redis.Client) *Service {
	return &Service{
		Repository: repository,
		Redis:      redis,
	}
}

// ===========================================

func (s *Service) ValidateUser(user *models.User) error {
	if len(user.Username) > 20 || len(user.Username) < 3 {
		logger.Error.Println(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if len(user.Password) > 20 || len(user.Password) < 6 {
		logger.Error.Println(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if strings.Contains(user.Password, "_") || strings.Contains(user.Password, "-") {
		logger.Error.Println(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if strings.Contains(user.Password, "@") || strings.Contains(user.Password, "#") {
		logger.Error.Println(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if strings.Contains(user.Password, "$") || strings.Contains(user.Password, "%") {
		logger.Error.Println(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if strings.Contains(user.Password, "&") || strings.Contains(user.Password, "*") {
		logger.Error.Println(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if strings.Contains(user.Password, "(") || strings.Contains(user.Password, ")") {
		logger.Error.Println(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if strings.Contains(user.Password, ":") || strings.Contains(user.Password, ".") {
		logger.Error.Println(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if strings.Contains(user.Password, "/") || strings.Contains(user.Password, `\`) {
		logger.Error.Println(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if strings.Contains(user.Password, ",") || strings.Contains(user.Password, ";") {
		logger.Error.Println(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if strings.Contains(user.Password, "?") || strings.Contains(user.Password, `"`) {
		logger.Error.Println(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	if strings.Contains(user.Password, "!") || strings.Contains(user.Password, "~") {
		logger.Error.Println(apperror.ErrForbidden)
		return apperror.ErrForbidden
	}
	return nil
}

func (s *Service) ExistsUser(username string) (bool, error) {
	existsUser, err := s.Repository.ExistsUser(username)
	if err != nil {
		logger.Error.Println(err)
		return false, err
	}
	if existsUser {
		logger.Info.Println("this username is registered")
		return true, nil
	}

	return false, nil
}

func (s *Service) CreateUser(ctx context.Context, user *models.User) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error.Println("failed to generate hash from password due error: ", err)
		return "", err
	}

	user.Password = string(hashPassword)

	userID, err := s.Repository.CreateUser(ctx, user)
	if err != nil {
		logger.Error.Println("failed to create user")
		return "", err
	}

	return userID, nil
}

func (s *Service) CheckUser(user *models.User) (string, error) {
	u, err := s.Repository.CheckUser(user)
	if err != nil {
		logger.Error.Println(err)
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		logger.Error.Println(err)
		return "", err
	}

	return u.ID, nil
}

func (s *Service) CreateToken(userID string) (*models.TokenDetails, error) {
	td := &models.TokenDetails{}

	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	// Creating Access Token
	err := os.Setenv("ACCESS_SECRET", "secret")
	if err != nil {
		return nil, err
	}

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	// Creating Refresh Token
	err = os.Setenv("REFRESH_SECRET", "secret")
	if err != nil {
		return nil, err
	}
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func (s *Service) CreateAuth(userID string, td *models.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) // converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := s.Redis.Set(td.AccessUuid, userID, at.Sub(now)).Err()
	if errAccess != nil {
		logger.Error.Println(errAccess)
		return errAccess
	}
	errRefresh := s.Redis.Set(td.RefreshUuid, userID, rt.Sub(now)).Err()
	if errRefresh != nil {
		logger.Error.Println(errRefresh)
		return errRefresh
	}

	return nil
}

func (s *Service) ExistsAccount(number string) (bool, error) {
	existsAccount, err := s.Repository.ExistsAccount(number)
	if err != nil {
		logger.Error.Println(err)
		return false, err
	}

	return existsAccount, nil
}

func (s *Service) CreateAccount(account *models.Account) error {
	err := s.Repository.CreateAccount(account)
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	return nil
}

func (s *Service) GetAccounts(userID string) ([]models.Account, error) {
	accounts, err := s.Repository.GetAccounts(userID)
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	return accounts, nil
}

func (s *Service) GetAccountById(userID, id string) (models.Account, error) {
	account, err := s.Repository.GetAccountById(userID, id)
	if err != nil {
		logger.Error.Println(err)
		return models.Account{}, err
	}

	return account, nil
}

func (s *Service) UpdateAccount(account *models.Account) error {
	err := s.Repository.UpdateAccount(account)
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	return nil
}

func (s *Service) CreateTransaction(tr *models.Transaction) error {
	err := s.Repository.CreateTransaction(tr)
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	return nil
}

func (s *Service) GetTransactions(accountID string) ([]models.Transaction, error) {
	tr, err := s.Repository.GetTransactions(accountID)
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	return tr, nil
}

func (s *Service) GetTransactionById(id string) (models.Transaction, error) {
	tr, err := s.Repository.GetTransactionById(id)
	if err != nil {
		logger.Error.Println(err)
		return models.Transaction{}, err
	}

	return tr, nil
}

func (s *Service) GetAccountInfoById(accountID string) (*models.Account, error) {
	account, err := s.Repository.GetAccountInfoById(accountID)
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	return account, nil
}

func (s *Service) GetReports(userID string, report *models.Report) (*excelize.File, error) {
	var transactions []models.Transaction

	tr, err := s.Repository.GetReports(report)
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	excelReports, err := s.GetExcelReports(userID, tr)
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	if report == (&models.Report{}) {
		accounts, err := s.GetAccounts(userID)
		if err != nil {
			logger.Error.Println(err)
			return nil, err
		}

		for _, account := range accounts {
			transaction, err := s.GetTransactions(account.ID)
			if err != nil {
				logger.Error.Println(err)
				return nil, err
			}
			transactions = append(transactions, transaction...)
		}
	}

	return excelReports, nil
}

func (s *Service) GetExcelReports(userID string, tr []models.Transaction) (*excelize.File, error) {
	excelFile := excelize.NewFile()

	sheet, err := excelFile.NewSheet("Отчёт")
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	err = excelFile.SetCellValue("Отчёт", "A1", "Имя пользователя")
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	err = excelFile.SetCellValue("Отчёт", "B1", "Название счёта")
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	err = excelFile.SetCellValue("Отчёт", "C1", "Тип операции")
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	err = excelFile.SetCellValue("Отчёт", "D1", "Сумма")
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	err = excelFile.SetCellValue("Отчёт", "E1", "Дата совершения операции")
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	u, err := s.GetUserInfoById(userID)
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}
	for i, transaction := range tr {
		i += 2

		account, err := s.GetAccountInfoById(transaction.AccountID)
		if err != nil {
			logger.Error.Println(err)
			return nil, err
		}

		err = excelFile.SetCellValue("Отчёт", "A"+strconv.Itoa(i), u.Username)
		if err != nil {
			logger.Error.Println(err)
			return nil, err
		}

		err = excelFile.SetCellValue("Отчёт", "B"+strconv.Itoa(i), account.Number)
		if err != nil {
			logger.Error.Println(err)
			return nil, err
		}

		err = excelFile.SetCellValue("Отчёт", "C"+strconv.Itoa(i), transaction.Type)
		if err != nil {
			logger.Error.Println(err)
			return nil, err
		}

		err = excelFile.SetCellValue("Отчёт", "D"+strconv.Itoa(i), transaction.Amount)
		if err != nil {
			logger.Error.Println(err)
			return nil, err
		}

		err = excelFile.SetCellValue("Отчёт", "E"+strconv.Itoa(i), transaction.CreatedAt)
		if err != nil {
			logger.Error.Println(err)
			return nil, err
		}
	}
	excelFile.SetActiveSheet(sheet)

	err = excelFile.SaveAs("report.xlsx")
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	return excelFile, nil
}

func (s *Service) GetUserInfoById(userID string) (*models.User, error) {
	u, err := s.Repository.GetUserInfoById(userID)
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	return u, nil
}
