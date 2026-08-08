package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Token struct {
	secret []byte
}

type Params struct {
	Secret string
}

func New(p Params) *Token {
	return &Token{secret: []byte(p.Secret)}
}

type claims struct {
	jwt.RegisteredClaims
}

type Info struct {
	UserID    uuid.UUID
	TokenID   string
	ExpiresAt time.Time
}

func (t *Token) BuildAccessToken(userID uuid.UUID) (string, error) {
	now := time.Now().UTC()

	claims := claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(2 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(t.secret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signedToken, nil
}

func (t *Token) ParseAccessToken(tokenString string) (Info, error) {
	var claims claims

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (any, error) {
			return []byte(t.secret), nil
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

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return Info{}, fmt.Errorf("%w: invalid subject", ErrInvalidToken)
	}

	if claims.ID == "" {
		return Info{}, fmt.Errorf("%w: token id is empty", ErrInvalidToken)
	}

	if claims.ExpiresAt == nil {
		return Info{}, fmt.Errorf("%w: expiration is empty", ErrInvalidToken)
	}

	return Info{
		UserID:    userID,
		TokenID:   claims.ID,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}
