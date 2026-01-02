package repo

import (
	"card-service/entity"
	"context"
)

type CardRepo interface {
	AddCard(ctx context.Context, card *entity.Card) error
	GetCardByNumber(ctx context.Context, cardNumber string) bool
	BlockCard(ctx context.Context, cardNumber string) error
	UnBlockCard(ctx context.Context, cardNumber string) error
	GetCardInfo(ctx context.Context, cardNumber string) (*entity.Card, error)
	Accrual(ctx context.Context, cardNumber string, amount int64) error
	WithdrawalMoney(ctx context.Context, cardNumber string, amount int64) error
}
