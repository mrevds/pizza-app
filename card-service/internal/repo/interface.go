package repo

import (
	"card-service/entity"
	"context"
)

type CardRepo interface {
	AddCard(ctx context.Context, card *entity.Card) error
	GetCardByNumber(ctx context.Context, cardNumber string) bool
}
