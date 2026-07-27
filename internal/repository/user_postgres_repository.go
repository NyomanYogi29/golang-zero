package repository

import (
	"context"
	"latihan/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, name, email, password, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
	)

	return err
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, name, email, password, created_at
		FROM users
		WHERE email = $1
	`

	var user model.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
