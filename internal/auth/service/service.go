package service

import (
	"context"
	"errors"
	"fmt"
	"project/internal/auth/domain"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type store interface {
	CreateUser(ctx context.Context, user domain.User) error
	GetUser(ctx context.Context, email string) (user domain.User, err error)
}

type blacklist interface {
	Revoke(ctx context.Context, tokenID string, ttl time.Duration) error
}

type tokenBuilder interface {
	BuildAccessToken(userID uuid.UUID) (string, error)
}

type Service struct {
	store     store
	blacklist blacklist
	token     tokenBuilder
}

type Params struct {
	Store     store
	Token     tokenBuilder
	Blacklist blacklist
}

func New(p Params) *Service {
	return &Service{
		store:     p.Store,
		token:     p.Token,
		blacklist: p.Blacklist,
	}
}

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

type LoginParams struct {
	Email    string
	Password string
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
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

		return LoginResult{}, fmt.Errorf("get user: %w", err)
	}

	if err := comparePassword(user.PasswordHash, req.Password); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return LoginResult{}, ErrInvalidCredentials
		}

		return LoginResult{}, fmt.Errorf("compare password: %w", err)
	}

	accessToken, err := s.token.BuildAccessToken(user.ID)
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		AccessToken: accessToken,
	}, nil
}

type LogoutParams struct {
	TokenID   string
	ExpiresAt time.Time
}

func (s *Service) Logout(ctx context.Context, p LogoutParams) error {
	ttl := time.Until(p.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("invalid token")
	}

	if err := s.blacklist.Revoke(ctx, p.TokenID, ttl); err != nil {
		return fmt.Errorf("revoke access token: %w", err)
	}
	return nil
}
