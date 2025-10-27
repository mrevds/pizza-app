package pg

import (
	"context"

	"user-service/client"
	"user-service/internal/entity"
	"user-service/internal/repository"
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

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	row := r.db.Pool.QueryRow(ctx, `
  SELECT id, first_name, phone_number, password, created_at, updated_at
  FROM users WHERE phone_number = $1
 `, username)

	var u entity.User
	if err := row.Scan(&u.ID, &u.FirstName, &u.PhoneNumber, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}
