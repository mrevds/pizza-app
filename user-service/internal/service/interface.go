package service

import (
	"context"
	"user-service/internal/entity"
)

type UserService interface {
	Register(ctx context.Context, input RegisterInput) (*entity.User, error)
}
