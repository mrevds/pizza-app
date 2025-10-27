package repository

import "user-service/internal/entity"

type UserRepository interface {
	Register(ctx context.Context, user *entity.User) error
}