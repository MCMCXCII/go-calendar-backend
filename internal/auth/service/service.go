package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"project/internal/auth/domain"
)

//go:generate mockgen -source=service.go -destination=./mocks/mocks.go -package=mocks

type Store interface {
	CreateUser(ctx context.Context, user domain.User) error
	GetUser(ctx context.Context, email string) (user domain.User, err error)
}

type Blacklist interface {
	Revoke(ctx context.Context, tokenID string, ttl time.Duration) error
}

type TokenBuilder interface {
	BuildAccessToken(userID uuid.UUID, expiration time.Duration) (string, error)
}

type Service struct {
	store           Store
	blacklist       Blacklist
	token           TokenBuilder
	tokenExpiration time.Duration
}

type Params struct {
	Store           Store
	Token           TokenBuilder
	Blacklist       Blacklist
	TokenExpiration time.Duration
}

func New(p Params) *Service {
	return &Service{
		store:           p.Store,
		token:           p.Token,
		blacklist:       p.Blacklist,
		tokenExpiration: p.TokenExpiration,
	}
}
