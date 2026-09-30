package auth

import "time"

type Account struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Session struct {
	ID        string    `json:"id"`
	AccountID string    `json:"accountId"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type PublicAccount struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (a Account) Public() PublicAccount {
	return PublicAccount{ID: a.ID, Username: a.Username}
}
