package repository

import (
	"context"
	"latihan/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}
