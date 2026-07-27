package storage

import (
	"context"
	"errors"
	"time"

	"github.com/gym-trainer/auth-service/internal/model"
	"github.com/jackc/pgx/v5"
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

func (s *TokenStorage) GetRefreshToken(
	ctx context.Context,
	token string,
) (*model.RefreshToken, error) {
	query := `
		SELECT token, user_id, expires_at
		FROM refresh_tokens
		WHERE token = $1
	`

	var refreshToken model.RefreshToken

	err := s.db.QueryRow(ctx, query, token).Scan(
		&refreshToken.Token,
		&refreshToken.UserID,
		&refreshToken.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &refreshToken, nil
}

func (s *TokenStorage) DeleteRefreshToken(
	ctx context.Context,
	token string,
) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE token = $1
	`

	_, err := s.db.Exec(ctx, query, token)
	if err != nil {
		return err
	}

	return nil
}
