package service_test

import (
	"context"
	"errors"
	"latihan/internal/model"
	"latihan/internal/service"
	"testing"
)

type MockUserRepository struct {
	users map[string]*model.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{users: make(map[string]*model.User)}
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *model.User) error {
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	if u, exists := m.users[email]; exists {
		return u, nil
	}

	return nil, errors.New("User not found")
}

func TestRegisterUser_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	userService := service.NewUserService(mockRepo)

	req := model.UserRegisterRequestSchema{
		Name:     "Yogi",
		Email:    "nyomanyogi@gmail.com",
		Password: "password123",
	}

	user, err := userService.Register(context.Background(), req)

	if err != nil {
		t.Fatalf("Ekspektasi no error, tapi dapat: %v", err)
	}

	if user.ID == "" {
		t.Errorf("Ekspektasi ID terisi UUID, tapi kosong")
	}

	if user.Password == "password123" {
		t.Errorf("Ekspektasi password di-hash, tapi masih plain-text")
	}
}

func TestRegisterUser_DuplicatedEmail(t *testing.T) {
	mockRepo := NewMockUserRepository()
	userService := service.NewUserService(mockRepo)

	req := model.UserRegisterRequestSchema{
		Name:     "Yogi",
		Email:    "nyomanyogi@gmail.com",
		Password: "password123",
	}

	_, _ = userService.Register(context.Background(), req)

	_, err := userService.Register(context.Background(), req)

	if err == nil {
		t.Errorf("Ekspektasi error saat email duplikat, tapi err nil")
	}
}
