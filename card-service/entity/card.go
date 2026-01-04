package entity

import "time"

type Card struct {
	ID         string    `json:"id" db:"id"`
	CardNumber string    `json:"card_number" db:"card_number"`
	DateExpire time.Time `json:"date_expire" db:"date_expire"`
	Balance    int64     `json:"balance" db:"balance"`
	Cvv        string    `json:"cvv" db:"cvv"`
	Currency   string    `json:"currency" db:"currency"`
	OwnerID    string    `json:"owner_id" db:"owner_id"`
	IsActive   bool      `json:"is_active" db:"is_active"`
}
