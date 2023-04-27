package repository

import (
	"context"
	"github.com/k4zb3k/project/internal/apperror"
	"github.com/k4zb3k/project/internal/models"
	"github.com/k4zb3k/project/pkg/logger"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	Connection *gorm.DB
}

func NewRepository(conn *gorm.DB) *Repository {
	return &Repository{
		Connection: conn,
	}
}

//==================================================

func (r *Repository) ExistsUser(username string) (bool, error) {
	var u models.User
	err := r.Connection.Where("username = ?", username).Error
	if err != nil {
		return false, err
	}
	if u == (models.User{}) {
		return false, nil
	}

	return true, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *models.User) (userID string, err error) {
	err = r.Connection.WithContext(ctx).Create(&user).Error
	if err != nil {
		logger.Error.Println("failed to create user")
		return "", err
	}

	return user.ID, nil
}

func (r *Repository) CheckUser(user *models.User) (u *models.User, err error) {
	if tx := r.Connection.Where("username = ?", user.Username).Find(&u); tx.Error != nil {
		logger.Error.Println("failed to user", err)
		return u, tx.Error
	}
	if u == (&models.User{}) {
		logger.Error.Println(apperror.ErrNotFound)
		return u, apperror.ErrNotFound
	}

	return u, nil
}

func (r *Repository) ExistsAccount(number string) (bool, error) {
	var acc *models.Account

	err := r.Connection.Where("number = ?", number).Find(&acc).Error
	if err != nil {
		logger.Error.Println(err)
		return false, err
	}
	if acc.Number != "" {
		return true, nil
	}

	return false, nil
}

func (r *Repository) CreateAccount(account *models.Account) error {
	err := r.Connection.Omit("created_at", "updated_at", "deleted_at").Create(&account).Error
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	return nil
}

func (r *Repository) GetAccounts(userID string) (accounts []models.Account, err error) {
	err = r.Connection.Where("user_id = ?", userID).Find(&accounts).Error
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	return accounts, nil
}

func (r *Repository) GetAccountById(userID, id string) (account models.Account, err error) {
	err = r.Connection.Where("user_id = ? and id = ?", userID, id).Find(&account).Error
	if err != nil {
		logger.Error.Println(err)
		return models.Account{}, err
	}

	return account, nil
}

func (r *Repository) UpdateAccount(account *models.Account) error {
	err := r.Connection.Save(account).Error
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	return nil
}

func (r *Repository) CreateTransaction(tr *models.Transaction) error {
	err := r.Connection.Omit("created_at", "updated_at", "deleted_at").Create(&tr).Error
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	return nil
}

func (r *Repository) GetTransactions(accountID string) (tr []models.Transaction, err error) {
	err = r.Connection.Where("account_id = ?", accountID).Find(&tr).Error
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	return tr, nil
}

func (r *Repository) GetTransactionById(id string) (tr models.Transaction, err error) {
	err = r.Connection.Where("id = ?", id).Find(&tr).Error
	if err != nil {
		logger.Error.Println(err)
		return models.Transaction{}, err
	}

	return tr, err
}

func (r *Repository) GetAccountInfoById(accountID string) (acc *models.Account, err error) {
	err = r.Connection.Where("id = ?", accountID).Find(&acc).Error
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	return acc, nil
}

func (r *Repository) GetReports(report *models.Report) (tr []models.Transaction, err error) {
	query := r.Connection
	if report.Type != "" {
		query = query.Where("type = ?", report.Type)
	}
	if report.From != (time.Time{}) {
		query = query.Where("created_at >= ?", report.From)
	}
	if report.To != (time.Time{}) {
		query = query.Where("created_at <= ?", report.To)
	}

	page := 1
	limit := 0

	if report.Page > 0 {
		page = report.Page
	}
	if report.Limit > 0 {
		limit = report.Limit
	}
	if report.Page > 0 {
		query = query.Limit(limit).Offset((page - 1) * limit)
	}

	err = query.Find(&tr).Error
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	return tr, nil
}

func (r *Repository) GetUserInfoById(userID string) (u *models.User, err error) {
	err = r.Connection.Where("id = ?", userID).Find(&u).Error
	if err != nil {
		logger.Error.Println(err)
		return nil, err
	}

	return u, nil
}
