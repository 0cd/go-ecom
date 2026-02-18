package auth

import (
	"context"
	"fmt"

	repo "github.com/0cd/go-ecom/internal/adapters/sqlc"
	"github.com/0cd/go-ecom/internal/users"
	"github.com/0cd/go-ecom/internal/utils"
)

type Service interface {
	Register(ctx context.Context, email, password string) (repo.User, error)
	Login(ctx context.Context, email, password string) (repo.FindUserByEmailRow, error)
	GenerateTokens(userID int64) (string, string, error)
}

type service struct {
	UserService users.Service
}

func NewService(userService users.Service) Service {
	return &service{UserService: userService}
}

func (s *service) Register(ctx context.Context, email, password string) (repo.User, error) {
	return s.UserService.CreateUser(ctx, users.CreateUserParams{
		Email:    email,
		Password: password,
	})
}

func (s *service) Login(ctx context.Context, email, password string) (repo.FindUserByEmailRow, error) {
	user, err := s.UserService.FindUserByEmail(ctx, email)
	if err != nil {
		return repo.FindUserByEmailRow{}, fmt.Errorf("invalid credentials")
	}
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return repo.FindUserByEmailRow{}, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

func (s *service) GenerateTokens(userID int64) (string, string, error) {
	accessToken, err := GenerateAccessToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := GenerateRefreshToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}
