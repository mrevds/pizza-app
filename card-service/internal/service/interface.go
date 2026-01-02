package service

import (
	"card-service/entity"
	"context"
)

type CardService interface {
	AddCard(ctx context.Context, input AddCardData) (*entity.Card, error)
	BlockCard(ctx context.Context, cardNumber string, userID string) error
	UnBlockCard(ctx context.Context, cardNumber string, userID string) error
	MoneyTransfer(ctx context.Context, userFrom string, userTo string, amount int64, userID string) error
}
