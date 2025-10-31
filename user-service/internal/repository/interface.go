package repository

import (
	"context"
	"user-service/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByPhoneNumber(ctx context.Context, username string) (*entity.User, error)
	SaveRefreshToken(ctx context.Context, rt *entity.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*entity.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RevokeUserRefreshTokens(ctx context.Context, userID string) error
	GetProfileInfo(ctx context.Context, userID string) (*entity.User, error)
	UpdateProfile(ctx context.Context, user *entity.User) error
}
