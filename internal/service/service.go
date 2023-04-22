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
	"golang.org/x/crypto/bcrypt"
	"os"
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
