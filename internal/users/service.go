package users

import (
	"context"
	"fmt"

	repo "github.com/0cd/go-ecom/internal/adapters/sqlc"
	"github.com/0cd/go-ecom/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	CreateUser(ctx context.Context, params CreateUserParams) (repo.User, error)
	DeleteUser(ctx context.Context, id int64) error

	ListUsers(ctx context.Context) ([]repo.ListUsersRow, error)
	SearchUsers(ctx context.Context, query pgtype.Text) ([]repo.SearchUsersRow, error)
	FindUserByID(ctx context.Context, id int64) (repo.FindUserByIDRow, error)
	FindUserByEmail(ctx context.Context, email string) (repo.FindUserByEmailRow, error)

	UpdateUser(ctx context.Context, updates repo.UpdateUserParams) (repo.User, error)
	UpdateUserPassword(ctx context.Context, id int64, oldPassword, newPassword string) error
	UpdateUserEmail(ctx context.Context, id int64, email string) error

	VerifyUser(ctx context.Context, id int64, verificationToken string) error
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &service{repo: repo}
}

func (s *service) CreateUser(ctx context.Context, params CreateUserParams) (repo.User, error) {
	if params.Email == "" {
		return repo.User{}, fmt.Errorf("email is required")
	}
	if params.Password == "" {
		return repo.User{}, fmt.Errorf("password is required")
	}

	if !utils.ValidateEmail(params.Email) {
		return repo.User{}, fmt.Errorf("invalid email")
	}
	if len(params.Password) < 8 {
		return repo.User{}, fmt.Errorf("password must be at least 8 characters long")
	}

	passwordHash, err := utils.HashPassword(params.Password)
	if err != nil {
		return repo.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	createdUser, err := s.repo.CreateUser(ctx, repo.CreateUserParams{
		Email:        params.Email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return repo.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}

func (s *service) DeleteUser(ctx context.Context, id int64) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *service) ListUsers(ctx context.Context) ([]repo.ListUsersRow, error) {
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		return []repo.ListUsersRow{}, fmt.Errorf("failed to fetch users: %w", err)
	}

	return users, nil
}

func (s *service) SearchUsers(ctx context.Context, query pgtype.Text) ([]repo.SearchUsersRow, error) {
	foundUsers, err := s.repo.SearchUsers(ctx, query)
	if err != nil {
		return []repo.SearchUsersRow{}, fmt.Errorf("failed to search users: %w", err)
	}

	return foundUsers, nil
}

func (s *service) FindUserByID(ctx context.Context, id int64) (repo.FindUserByIDRow, error) {
	foundUser, err := s.repo.FindUserByID(ctx, id)
	if err != nil {
		return repo.FindUserByIDRow{}, fmt.Errorf("failed to fetch user (id: %v): %w", id, err)
	}

	return foundUser, nil
}

func (s *service) FindUserByEmail(ctx context.Context, email string) (repo.FindUserByEmailRow, error) {
	foundUser, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return repo.FindUserByEmailRow{}, fmt.Errorf("failed to fetch user (email: %s): %w", email, err)
	}

	return foundUser, nil
}

func (s *service) UpdateUser(ctx context.Context, updates repo.UpdateUserParams) (repo.User, error) {
	_, err := s.repo.FindUserByID(ctx, updates.ID)
	if err != nil {
		return repo.User{}, fmt.Errorf("failed to find user %w", err)
	}

	if updates.Email.Valid {
		if updates.Email.String == "" {
			return repo.User{}, fmt.Errorf("email cannot be empty")
		}
		if !utils.ValidateEmail(updates.Email.String) {
			return repo.User{}, fmt.Errorf("invalid email")
		}
		// require email verification after update by default
		updates.Verified.Bool = false
	}

	// didn't feel like creating a separate struct with password variable
	// receive plain text password and then hash it
	if updates.PasswordHash.Valid {
		if updates.PasswordHash.String == "" {
			return repo.User{}, fmt.Errorf("password cannot be empty")
		}

		if len(updates.PasswordHash.String) < 8 {
			return repo.User{}, fmt.Errorf("password must be at least 8 characters long")
		}

		hashed, err := utils.HashPassword(updates.PasswordHash.String)
		if err != nil {
			return repo.User{}, fmt.Errorf("failed to hash password: %w", err)
		}

		updates.PasswordHash.String = hashed
	}

	if !updates.Email.Valid && !updates.PasswordHash.Valid && !updates.Verified.Valid {
		return repo.User{}, fmt.Errorf("no fields to update")
	}

	updatedUser, err := s.repo.UpdateUser(ctx, updates)
	if err != nil {
		return repo.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	return updatedUser, nil
}

func (s *service) UpdateUserPassword(ctx context.Context, id int64, oldPassword, newPassword string) error {
	user, err := s.repo.FindUserByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	if oldPassword == "" {
		return fmt.Errorf("current password is required")
	}
	if newPassword == "" {
		return fmt.Errorf("new password is required")
	}
	if len(newPassword) < 8 {
		return fmt.Errorf("new password must be at least 8 characters long")
	}

	if !utils.CheckPasswordHash(oldPassword, user.PasswordHash) {
		return fmt.Errorf("current password is incorrect")
	}

	newHashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	_, err = s.repo.UpdateUser(ctx, repo.UpdateUserParams{
		ID: id,
		PasswordHash: pgtype.Text{
			Valid:  true,
			String: newHashed,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (s *service) UpdateUserEmail(ctx context.Context, id int64, email string) error {
	_, err := s.repo.FindUserByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if !utils.ValidateEmail(email) {
		return fmt.Errorf("invalid email")
	}

	_, err = s.repo.UpdateUser(ctx, repo.UpdateUserParams{
		ID: id,
		Email: pgtype.Text{
			Valid:  true,
			String: email,
		},
		Verified: pgtype.Bool{
			Valid: true,
			Bool:  false,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update email: %w", err)
	}

	return nil
}

func (s *service) VerifyUser(ctx context.Context, id int64, verificationToken string) error {
	user, err := s.repo.FindUserByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	if user.Verified {
		return fmt.Errorf("user already verified")
	}

	// TODO: generate an email verification token at user registration and include it in the link
	// TODO: verify it against db value here
	// TODO: implement the feature lol
	// just a static string for now

	if verificationToken == "" {
		return fmt.Errorf("verification token is required")
	}

	if verificationToken != "forsenE" {
		return fmt.Errorf("invalid verification token")
	}

	_, err = s.repo.UpdateUser(ctx, repo.UpdateUserParams{
		ID: id,
		Verified: pgtype.Bool{
			Valid: true,
			Bool:  true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to verify user: %w", err)
	}

	return nil
}
