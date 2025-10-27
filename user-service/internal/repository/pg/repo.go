package pg

import (
	"user-service/internal/repository"
	"user-service/client")

type UserRepo struct {
	db *client.DB
}

func NewUserRepo(db *client.DB) repository.UserRepository {
	return &UserRepo{db: db}
}