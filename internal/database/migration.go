package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var allSchemas = []string{
	UserSchema,
}

func RunMigrations(ctx context.Context, db *pgxpool.Pool) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed to starting migration: %w", err)
	}

	defer tx.Rollback(ctx)

	log.Println("Starting database migrations...")

	for idx, schema := range allSchemas {
		if _, err := tx.Exec(ctx, schema); err != nil {
			return fmt.Errorf("Failed migrating schema on index: %d: %w", idx+1, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("Failed to commit migration: %w", err)
	}

	log.Println("All database migrations completed successfully!")
	return nil
}
