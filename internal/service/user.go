package service

import (
	"context"
	"time"

	"github.com/gym-trainer/auth-service/internal/model"
	"github.com/gym-trainer/auth-service/internal/storage"
	"github.com/gym-trainer/auth-service/internal/token"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo             *storage.UserStorage
	tokenRepo            *storage.TokenStorage
	maker                *token.Maker
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewUserService(
	userRepo *storage.UserStorage,
	tokenRepo *storage.TokenStorage,
	maker *token.Maker,
) *UserService {
	return &UserService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		maker:     maker,
	}
}

func (s *UserService) issueTokens(
	ctx context.Context,
	userID int,
) (*model.AuthResult, error) {
	timeNow := time.Now()

	accessTokenExpiresAt := timeNow.Add(s.accessTokenDuration)
	accessToken, err := s.maker.CreateAccessToken(
		userID,
		timeNow,
		accessTokenExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.maker.CreateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshTokenExpiresAt := timeNow.Add(s.refreshTokenDuration)
	err = s.tokenRepo.StoreRefreshToken(
		ctx,
		refreshToken,
		userID,
		refreshTokenExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	return &model.AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) Register(
	ctx context.Context,
	input model.RegisterInput,
) (*model.AuthResult, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.CreateUser(ctx, input.Email, string(hashedBytes))
	if err != nil {
		return nil, err
	}

	result, err := s.issueTokens(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	result.User = user

	return result, nil
}

func (s *UserService) Login(
	ctx context.Context,
	input model.LoginInput,
) (*model.AuthResult, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, model.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, model.ErrInvalidCredentials
	}

	result, err := s.issueTokens(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	result.User = user

	return result, nil
}

func (s *UserService) Refresh(
	ctx context.Context,
	refreshToken string,
) (*model.AuthResult, error) {
	token, err := s.tokenRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, model.ErrUnauthorized
	}

	timeNow := time.Now()
	s.tokenRepo.DeleteRefreshToken(ctx, refreshToken)
	if token.ExpiresAt.Before(timeNow) {
		return nil, model.ErrUnauthorized
	}

	return s.issueTokens(ctx, token.UserID)
}
