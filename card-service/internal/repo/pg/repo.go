package pg

import (
	"card-service/database"
	"card-service/entity"
	"card-service/internal/repo"
	"context"
)

type cardRepo struct {
	db *database.DB
}

func NewCardRepo(dbb *database.DB) repo.CardRepo {
	return &cardRepo{db: dbb}
}

func (r *cardRepo) AddCard(ctx context.Context, card *entity.Card) error {
	_, err := r.db.Pool.Exec(ctx, `
	INSERT INTO cards (id, card_number, date_expire, balance, currency, owner_id, is_active) VALUES ($1,$2,$3,$4,$5,$6,$7)`, card.ID, card.CardNumber, card.DateExpire, card.Balance, card.Currency, card.OwnerID, card.IsActive)
	return err
}

func (r *cardRepo) GetCardByNumber(ctx context.Context, cardNumber string) bool {
	var id int
	err := r.db.Pool.QueryRow(ctx, `
	SELECT id from cards where card_number = $1`, cardNumber).Scan(&id)
	if err != nil {
		return false
	}
	return true
}
