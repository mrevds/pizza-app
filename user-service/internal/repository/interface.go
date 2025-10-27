package repository

import (
	"context"
	"user-service/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
}
