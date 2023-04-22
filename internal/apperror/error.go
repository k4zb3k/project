package apperror

import (
	"encoding/json"
	"github.com/k4zb3k/project/pkg/logger"
)

var (
	ErrUnauthorized   = NewAppError(nil, "please provide valid login details", "", "US-000008")
	ErrBadRequest     = NewAppError(nil, "bad request data", "", "US-000007")
	ErrNotFound       = NewAppError(nil, "not found", "", "US-000003")
	ErrForbidden      = NewAppError(nil, "forbidden", "", "US-000001")
	ErrRegistered     = NewAppError(nil, "user already registered", "", "US-000002")
	ErrInvalidToken   = NewAppError(nil, "invalid jwt token", "", "US-000004")
	ErrInvalid        = NewAppError(nil, "validate error", "", "US-000005")
	ErrInternalServer = NewAppError(nil, "internal server error", "", "US-000006")
	ErrExpiredRefresh = NewAppError(nil, "refresh token is expired", "", "US-000008")
	ErrExpiredToken   = NewAppError(nil, "token is expired", "", "US-000009")
	ErrExistsAccount  = NewAppError(nil, "account is exists", "", "US-000010")
)

type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
	Code             string `json:"code,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Marshal() []byte {
	marshal, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		logger.Error.Fatalln(err)
		return nil
	}
	return marshal
}

func NewAppError(err error, message, developerMessage, code string) *AppError {
	return &AppError{
		Err:              err,
		Message:          message,
		DeveloperMessage: developerMessage,
		Code:             code,
	}
}

func systemError(err error) *AppError {
	return NewAppError(err, "internal system error", err.Error(), "US-000000")
}
