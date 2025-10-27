package service

import (
	"fmt"
	"time"
	"user-service/internal/entity"
	"user-service/internal/repository"

	"context"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Username string
	Email    string
	Password string
	Name     *string
	Age      *int32
	Bio      *string
}
type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(ctx context.Context, input RegisterInput) (*entity.User, error) {
	existing, _ := s.repo.GetByUsername(ctx, input.Username)
	if existing != nil {
		return nil, fmt.Errorf("username already taken")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID: uuid.NewString(),

		Email: input.Email,

		Password:  string(hashed),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
