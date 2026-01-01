package service

import (
	"card-service/entity"
	"card-service/internal/repo"
	"card-service/utils"
	"context"
	"fmt"
	"time"
)

type AddCardData struct {
	Id         string
	CardNumber string
	DateExpire time.Time
	Cvv        string
	Currency   string
	OwnerID    string
	IsActive   bool
}

type cardService struct {
	cardRepo repo.CardRepo
}

func NewCardService(cardRepo repo.CardRepo) CardService {
	return &cardService{
		cardRepo: cardRepo,
	}
}

func (s *cardService) AddCard(ctx context.Context, input AddCardData) (*entity.Card, error) {
	isExists := s.cardRepo.GetCardByNumber(ctx, input.CardNumber)
	if isExists {
		return nil, fmt.Errorf("card already exist")
	}

	hashedCVV := utils.CvvHash(input.Cvv)

	card := &entity.Card{
		ID:         input.Id,
		CardNumber: input.CardNumber,
		DateExpire: input.DateExpire,
		Cvv:        hashedCVV,
		Currency:   input.Currency,
		OwnerID:    input.OwnerID,
		IsActive:   input.IsActive,
	}
	err := s.cardRepo.AddCard(ctx, card)
	if err != nil {
		fmt.Printf("AddCard error: %v\n", err)
		return nil, fmt.Errorf("failed to add card: %w", err)
	}
	return card, nil
}
