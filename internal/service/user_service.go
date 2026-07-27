package service

import (
	"context"
	"errors"
	"fmt"
	"latihan/internal/model"
	"latihan/internal/repository"
	"latihan/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
	rdb  *redis.Client
}

func NewUserService(repo repository.UserRepository, rdb *redis.Client) *UserService {
	return &UserService{
		repo: repo,
		rdb:  rdb,
	}
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

func (s *UserService) Login(ctx context.Context, req model.UserLoginRequestSchema) (*model.UserLoginResponseSchema, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("Invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("Invalid email or password")
	}

	tokenString, duration, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("gagal generate token: %w", err)
	}

	sessionKey := fmt.Sprintf("session:%s", user.ID)
	err = s.rdb.Set(ctx, sessionKey, tokenString, duration).Err()
	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan session di redis: %w", err)
	}

	return &model.UserLoginResponseSchema{
		Token:     tokenString,
		TokenType: "Bearer",
	}, nil
}

func (s *UserService) IsSessionValid(ctx context.Context, userId, tokenString string) bool {
	sessionKey := fmt.Sprintf("session:%s", userId)
	savedToken, err := s.rdb.Get(ctx, sessionKey).Result()
	if err != nil {
		return false
	}

	return savedToken == tokenString
}

func (s *UserService) Logout(ctx context.Context, userID string) error {
	sessionKey := fmt.Sprintf("session:%s", userID)
	return s.rdb.Del(ctx, sessionKey).Err()
}
