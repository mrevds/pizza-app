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

func (s *cardService) BlockCard(ctx context.Context, cardNumber string, userID string) error {
	if len(cardNumber) != 16 {
		return fmt.Errorf("len not equal 16")
	}
	owner, err := s.cardRepo.GetCardInfo(ctx, cardNumber)
	if err != nil {
		return err
	}
	if owner.OwnerID != userID {
		return fmt.Errorf("ne owner blyad")
	}
	err = s.cardRepo.BlockCard(ctx, cardNumber)
	if err != nil {
		return err
	}
	return nil
}

func (s *cardService) UnBlockCard(ctx context.Context, cardNumber string, userID string) error {
	if len(cardNumber) != 16 {
		return fmt.Errorf("len not equal 16")
	}
	owner, err := s.cardRepo.GetCardInfo(ctx, cardNumber)
	if err != nil {
		return err
	}
	if owner.OwnerID != userID {
		return fmt.Errorf("ne owner blyad")
	}
	err = s.cardRepo.UnBlockCard(ctx, cardNumber)
	if err != nil {
		return err
	}
	return nil
}

func (s *cardService) MoneyTransfer(ctx context.Context, userFrom string, userTo string, amount int64, userID string) error {
	if len(userFrom) != 16 || len(userTo) != 16 {
		return fmt.Errorf("len not enough")
	}
	cardOwner, _ := s.cardRepo.GetCardInfo(ctx, userFrom)
	if cardOwner.OwnerID != userID {
		return fmt.Errorf("cardowner is wrong:%s", cardOwner.OwnerID)
	}
	exist := s.cardRepo.GetCardByNumber(ctx, userTo)
	if !exist {
		return fmt.Errorf("cant find card: %t", exist)
	}

	err := s.cardRepo.WithdrawalMoney(ctx, userFrom, amount)
	if err != nil {
		return fmt.Errorf("withdraw err: %w", err)
	}
	err = s.cardRepo.Accrual(ctx, userTo, amount)
	if err != nil {
		return fmt.Errorf("accrual err: %w", err)
	}
	return nil
}
