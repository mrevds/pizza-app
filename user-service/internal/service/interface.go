package service

import (
	"context"
	"user-service/internal/entity"
)

type UserService interface {
	Register(ctx context.Context, input RegisterInput) (*entity.User, error)
	Login(ctx context.Context, phoneNumber, password string) (user *entity.User, accessToken, refreshToken string)
	RefreshTokens(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, err error)
	GetProfileInfo(ctx context.Context, userID string) (*entity.User, error)
}
