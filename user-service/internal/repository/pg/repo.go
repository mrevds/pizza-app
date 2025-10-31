package pg

import (
	"context"
	"fmt"
	"user-service/client"
	"user-service/internal/entity"
	"user-service/internal/repository"

	"github.com/jackc/pgx/v4"
)

type userRepo struct {
	db *client.DB
}

func NewUserRepo(db *client.DB) repository.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *entity.User) error {
	_, err := r.db.Pool.Exec(ctx, `
  INSERT INTO users (id, first_name, phone_number, password, created_at, updated_at)
  VALUES ($1, $2, $3, $4, $5, $6)
 `, u.ID, u.FirstName, u.PhoneNumber, u.Password, u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *userRepo) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.User, error) {
	row := r.db.Pool.QueryRow(ctx, `
	  SELECT id, first_name, phone_number, password, created_at, updated_at FROM users WHERE phone_number = $1`, phoneNumber)
	var u entity.User
	if err := row.Scan(&u.ID, &u.FirstName, &u.PhoneNumber, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) SaveRefreshToken(ctx context.Context, rt *entity.RefreshToken) error {
	_, err := r.db.Pool.Exec(ctx, `
  INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at, revoked)
  VALUES ($1, $2, $3, $4, $5, $6)
 `, rt.ID, rt.UserID, rt.Token, rt.ExpiresAt, rt.CreatedAt, rt.Revoked)
	return err
}

func (r *userRepo) GetRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error) {
	var rt entity.RefreshToken
	err := r.db.Pool.QueryRow(ctx, `SELECT id, user_id, token, expires_at, created_at, revoked FROM refresh_tokens WHERE token = $1 AND revoked = false`, token).Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt, &rt.Revoked)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &rt, nil
}

func (r *userRepo) RevokeRefreshToken(ctx context.Context, token string) error {
	_, err := r.db.Pool.Exec(ctx, `
        UPDATE refresh_tokens SET revoked = true WHERE token = $1
    `, token)
	return err
}

func (r *userRepo) RevokeUserRefreshTokens(ctx context.Context, userID string) error {
	_, err := r.db.Pool.Exec(ctx, `
        UPDATE refresh_tokens SET revoked = true WHERE user_id = $1
    `, userID)
	return err
}

func (r *userRepo) GetProfileInfo(ctx context.Context, userID string) (*entity.User, error) {
	row := r.db.Pool.QueryRow(ctx, `
	  SELECT id, first_name, phone_number, password, created_at, updated_at FROM users WHERE id = $1`, userID)
	var u entity.User
	if err := row.Scan(&u.ID, &u.FirstName, &u.PhoneNumber, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) UpdateProfile(ctx context.Context, user *entity.User) error {
	query := "UPDATE users SET updated_at = $1"
	args := []interface{}{user.UpdatedAt}
	argNum := 2

	if user.FirstName != "" {
		query += fmt.Sprintf(", first_name = $%d", argNum)
		args = append(args, user.FirstName)
		argNum++
	}
	if user.LastName != "" {
		query += fmt.Sprintf(", last_name = $%d", argNum)
		args = append(args, user.LastName)
		argNum++
	}
	if user.Email != "" {
		query += fmt.Sprintf(", email = $%d", argNum)
		args = append(args, user.Email)
		argNum++
	}
	if user.PhoneNumber != "" {
		query += fmt.Sprintf(", phone_number = $%d", argNum)
		args = append(args, user.PhoneNumber)
		argNum++
	}

	query += fmt.Sprintf(" WHERE id = $%d", argNum)
	args = append(args, user.ID)

	_, err := r.db.Pool.Exec(ctx, query, args...)
	return err
}
