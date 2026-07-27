package token

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Maker struct {
	privateKey *rsa.PrivateKey
}

func NewMaker(privateKeyPEM []byte) (*Maker, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &Maker{privateKey: privateKey}, nil
}

func (m *Maker) CreateAccessToken(
	userID int,
	createdAt time.Time,
	expiresAt time.Time,
) (string, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate token ID: %w", err)
	}

	claims := jwt.MapClaims{
		"sub": userID,
		"jti": tokenID,
		"iat": createdAt.Unix(),
		"exp": expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(m.privateKey)
}

func (m *Maker) CreateRefreshToken() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func (m *Maker) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return &m.privateKey.PublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
