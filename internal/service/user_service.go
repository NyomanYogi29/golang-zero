package service

import (
	"context"
	"errors"
	"latihan/internal/model"
	"latihan/internal/repository"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, req model.UserRegisterRequestSchema) (*model.User, error) {
	existingUser, _ := s.repo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("Email sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &model.User{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateUser(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
