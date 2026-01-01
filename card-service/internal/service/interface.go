package service

import (
	"card-service/entity"
	"context"
)

type CardService interface {
	AddCard(ctx context.Context, input AddCardData) (*entity.Card, error)
}