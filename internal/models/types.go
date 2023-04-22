package models

import "time"

type User struct {
	ID       string `gorm:"type:uuid;default:uuid_generate_v4()"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessUuid   string `json:"access_uuid"`
	RefreshUuid  string `json:"refresh_uuid"`
	AtExpires    int64  `json:"at_expires"`
	RtExpires    int64  `json:"rt_expires"`
}

type AccessDetails struct {
	AccessUuid string `json:"access_uuid"`
	UserId     string `json:"user_id"`
}

type Token struct {
	ID    string `json:"-"`
	Token string `json:"token"`
}

type Account struct {
	ID        string    `gorm:"type:uuid;default:uuid_generate_v4()"`
	UserID    string    `json:"user_id,omitempty"`
	Number    string    `json:"number"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

type Transaction struct {
	ID        string  `gorm:"type:uuid;default:uuid_generate_v4()"`
	AccountID string  `json:"account_id"`
	Type      string  `json:"type"`
	Amount    float64 `json:"amount"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type Report struct {
	ID        string    `gorm:"type:uuid;default:uuid_generate_v4()"`
	AccountID string    `json:"account_id,omitempty"`
	Type      string    `json:"type,omitempty"`
	Limit     int       `json:"limit,omitempty"`
	Page      int       `json:"page,omitempty"`
	DateFrom  string    `json:"date_from,omitempty"`
	DateTo    string    `json:"date_to,omitempty"`
	From      time.Time `json:"-"`
	To        time.Time `json:"-"`
}
