package pg

import (
	"card-service/database"
	"card-service/entity"
	"card-service/internal/repo"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
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
	var id string
	err := r.db.Pool.QueryRow(ctx, `
	SELECT id from cards where card_number = $1`, cardNumber).Scan(&id)
	if err != nil {
		return false
	}
	return true
}

func (r *cardRepo) GetCardInfo(ctx context.Context, cardNumber string) (*entity.Card, error) {
	card := &entity.Card{}
	err := r.db.Pool.QueryRow(ctx, `SELECT id, card_number, date_expire, balance, currency, owner_id, is_active from cards where card_number = $1`, cardNumber).Scan(
		&card.ID,
		&card.CardNumber,
		&card.DateExpire,
		&card.Balance,
		&card.Currency,
		&card.OwnerID,
		&card.IsActive,
	)
	if err != nil {
		return nil, err
	}
	return card, nil
}

func (r *cardRepo) BlockCard(ctx context.Context, cardNumber string) error {
	_, err := r.db.Pool.Exec(ctx, `
		UPDATE cards SET is_active = false WHERE card_number = $1`, cardNumber)
	return err
}

func (r *cardRepo) UnBlockCard(ctx context.Context, cardNumber string) error {
	_, err := r.db.Pool.Exec(ctx, `
	UPDATE cards SET is_active = true WHERE card_number = $1`, cardNumber)
	return err
}

func (r *cardRepo) WithdrawalMoney(ctx context.Context, cardNumber string, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(
		ctx,
		`UPDATE cards
	SET balance = balance - $2
	WHERE card_number = $1 AND balance >= $2`,
		cardNumber,
		amount,
	)
	if err != nil {
		return fmt.Errorf("error withdraw exec: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("insufficient funds or card not found")
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}
	return nil

}

func (r *cardRepo) Accrual(ctx context.Context, cardNumber string, amount int64) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var newBalance int64
	err = tx.QueryRow(ctx, `
    UPDATE cards 
    SET balance = balance + $1 
    WHERE card_number = $2 
    RETURNING balance
`, amount, cardNumber).Scan(&newBalance)
	if err == pgx.ErrNoRows {
		return fmt.Errorf("card not found")
	}
	if err != nil {
		return fmt.Errorf("error accrual: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}
	return nil
}
