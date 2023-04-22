package repository

import (
	"context"
	"github.com/k4zb3k/project/internal/apperror"
	"github.com/k4zb3k/project/internal/models"
	"github.com/k4zb3k/project/pkg/logger"
	"gorm.io/gorm"
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
	err := r.Connection.Omit("created", "updated", "deleted").Create(&account).Error
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
