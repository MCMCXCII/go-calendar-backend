package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"project/internal/auth/domain"
)

type LoginParams struct {
	Email    string
	Password string
}

type LoginResult struct {
	AccessToken string
}

func (s *Service) Login(ctx context.Context, req LoginParams) (LoginResult, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if email == "" || req.Password == "" {
		return LoginResult{}, ErrInvalidCredentials
	}

	user, err := s.store.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return LoginResult{}, ErrInvalidCredentials
		}
		return LoginResult{}, fmt.Errorf("store.GetUser: %w", err)
	}

	if err = comparePassword(user.PasswordHash, req.Password); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return LoginResult{}, ErrInvalidCredentials
		}
		return LoginResult{}, fmt.Errorf("compare password: %w", err)
	}

	accessToken, err := s.token.BuildAccessToken(user.ID, s.tokenExpiration)
	if err != nil {
		return LoginResult{}, fmt.Errorf("token.BuildAccessToken: %w", err)
	}

	return LoginResult{AccessToken: accessToken}, nil
}

func comparePassword(passwordHash, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(password),
	)
}
