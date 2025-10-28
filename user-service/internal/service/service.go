package service

import (
	"fmt"
	"time"
	"user-service/internal/entity"
	"user-service/internal/repository"
	"user-service/internal/utils"

	"context"

	"github.com/google/uuid"
)

type RegisterInput struct {
	FirstName   string
	PhoneNumber string
	Password    string
}
type userService struct {
	repo       repository.UserRepository
	jwtManager *utils.JWTManager
}

func NewUserService(repo repository.UserRepository, jwtManager *utils.JWTManager) UserService {
	return &userService{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

func (s *userService) Register(ctx context.Context, input RegisterInput) (*entity.User, error) {
	existing, _ := s.repo.GetByPhoneNumber(ctx, input.PhoneNumber)
	if existing != nil {
		return nil, fmt.Errorf("phone already taken")
	}

	hashed := utils.PasswordHash(input.Password)

	user := &entity.User{
		ID:          uuid.NewString(),
		FirstName:   input.FirstName,
		PhoneNumber: input.PhoneNumber,
		Password:    hashed,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
func (s *userService) Login(ctx context.Context, phoneNumber, password string) (user *entity.User, accessToken, refreshToken string) {
	if phoneNumber == "" || password == "" {
		return nil, "", ""
	}
	existing, err := s.repo.GetByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return nil, "", ""
	}
	if !utils.CheckPasswordHash(password, existing.Password) {
		return nil, "", ""
	}
	acToken, err := s.jwtManager.GenerateToken(existing.ID)
	if err != nil {
		return nil, "", ""
	}
	rfToken, err := s.jwtManager.GenerateRefreshToken(existing.ID)
	if err != nil {
		return nil, "", ""
	}

	// Сохраняем refresh token в базу
	refreshTokenEntity := &entity.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    existing.ID,
		Token:     rfToken,
		ExpiresAt: time.Now().Add(time.Duration(s.jwtManager.RefreshTokenTTL()) * time.Minute),
		CreatedAt: time.Now(),
		Revoked:   false,
	}
	if err := s.repo.SaveRefreshToken(ctx, refreshTokenEntity); err != nil {
		return nil, "", ""
	}

	return existing, acToken, rfToken
}

func (s *userService) RefreshTokens(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	// Валидируем refresh token
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Проверяем что это именно refresh token
	if claims.Type != "refresh" {
		return "", "", fmt.Errorf("token is not a refresh token")
	}

	// Проверяем что токен есть в базе и не отозван
	storedToken, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil || storedToken == nil {
		return "", "", fmt.Errorf("refresh token not found or revoked")
	}

	// Отзываем старый refresh token
	if err := s.repo.RevokeRefreshToken(ctx, refreshToken); err != nil {
		return "", "", fmt.Errorf("failed to revoke old token: %w", err)
	}

	// Генерируем новую пару токенов
	newAccessToken, err = s.jwtManager.GenerateToken(claims.UserID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err = s.jwtManager.GenerateRefreshToken(claims.UserID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Сохраняем новый refresh token в базу
	refreshTokenEntity := &entity.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    claims.UserID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(time.Duration(s.jwtManager.RefreshTokenTTL()) * time.Minute),
		CreatedAt: time.Now(),
		Revoked:   false,
	}
	if err := s.repo.SaveRefreshToken(ctx, refreshTokenEntity); err != nil {
		return "", "", fmt.Errorf("failed to save refresh token: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}
