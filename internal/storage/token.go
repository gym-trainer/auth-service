package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenStorage struct {
	db *pgxpool.Pool
}

func NewTokenStorage(db *pgxpool.Pool) *TokenStorage {
	return &TokenStorage{
		db: db,
	}
}

func (s *TokenStorage) StoreRefreshToken(
	ctx context.Context,
	token string,
	userID int,
	expiresAt time.Time,
) error {
	query := `
		INSERT INTO refresh_tokens (token, user_id, expires_at)
		VALUES ($1, $2, $3)
	`

	_, err := s.db.Exec(ctx, query, token, userID, expiresAt)

	return err
}
