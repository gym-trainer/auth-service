package service

import (
	"context"

	"github.com/gym-trainer/auth-service/internal/model"
	"github.com/gym-trainer/auth-service/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *storage.UserStorage
}

func NewUserService(repo *storage.UserStorage) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Register(
	ctx context.Context,
	input model.RegisterInput,
) (*model.User, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return s.repo.CreateUser(ctx, input.Email, string(hashedBytes))
}

func (s *UserService) Login(
	ctx context.Context,
	input model.LoginInput,
) (*model.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, input.Email)
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

	return user, nil
}
