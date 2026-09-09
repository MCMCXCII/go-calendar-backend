package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"project/internal/auth/domain"
)

type RegisterParams struct {
	Email    string
	Password string
}

type RegisterResult struct {
	UserID uuid.UUID
}

func (s *Service) Register(ctx context.Context, req RegisterParams) (RegisterResult, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" {
		return RegisterResult{}, ErrEmailEmpty
	}
	if req.Password == "" {
		return RegisterResult{}, ErrPasswordEmpty
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return RegisterResult{}, fmt.Errorf("error hash password: %w", err)
	}

	user := domain.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
	}

	if err := s.store.CreateUser(ctx, user); err != nil {
		return RegisterResult{}, fmt.Errorf("error create user: %w", err)
	}

	return RegisterResult{
		UserID: user.ID,
	}, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", fmt.Errorf("generate password hash: %w", err)
	}

	return string(hash), nil
}
