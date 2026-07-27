package config

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

func ConnectPostgres() error {
	dsn := "postgres://postgres:password@localhost:5432/app_db?sslmode=disable"

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("gagal parse config postgres: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("gagal konek ke postgres: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("gagal ping postgres: %w", err)
	}

	dbPool = pool

	return nil
}

func GetDBPool() *pgxpool.Pool {
	return dbPool
}
