package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Config struct {
	SecretKey string `envconfig:"JWT_SECRET" required:"true"`
}

type Token struct {
	secret []byte
}

func New(c Config) *Token {
	return &Token{secret: []byte(c.SecretKey)}
}

type claims struct {
	jwt.RegisteredClaims
}

type Info struct {
	TokenID   string
	ExpiresAt time.Time
	UserID    uuid.UUID
}

func (t *Token) BuildAccessToken(userID uuid.UUID, expirations time.Duration) (string, error) {
	now := time.Now().UTC()

	c := claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expirations)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	signedToken, err := token.SignedString(t.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signedToken, nil
}

func (t *Token) ParseAccessToken(tokenString string) (Info, error) {
	var c claims

	token, err := jwt.ParseWithClaims(
		tokenString,
		&c,
		func(_ *jwt.Token) (any, error) {
			return t.secret, nil
		},
		jwt.WithValidMethods([]string{
			jwt.SigningMethodHS256.Alg(),
		}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return Info{}, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !token.Valid {
		return Info{}, ErrInvalidToken
	}

	userID, err := uuid.Parse(c.Subject)
	if err != nil {
		return Info{}, fmt.Errorf("%w: invalid subject", ErrInvalidToken)
	}

	if c.ID == "" {
		return Info{}, fmt.Errorf("%w: token id is empty", ErrInvalidToken)
	}

	if c.ExpiresAt == nil {
		return Info{}, fmt.Errorf("%w: expiration is empty", ErrInvalidToken)
	}

	return Info{
		UserID:    userID,
		TokenID:   c.ID,
		ExpiresAt: c.ExpiresAt.Time,
	}, nil
}
