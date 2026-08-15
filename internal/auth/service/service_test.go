package service

import (
	"context"
	"project/internal/auth/domain"
	"testing"

	"github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:generate mockgen -source=service.go -destination=mocks_test.go -package=service

func TestService_Register(t *testing.T) {
	ctx := context.Background()

	t.Run("успешная регистрация", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		store := NewMockstore(ctrl)

		var savedUser domain.User
		store.EXPECT().
			CreateUser(ctx, gomock.Any()).
			DoAndReturn(func(_ context.Context, user domain.User) error {
				savedUser = user
				return nil
			})

		svc := New(Params{Store: store})

		result, err := svc.Register(ctx, RegisterParams{
			Email:    "Test@Example.com",
			Password: "password123",
		})

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, result.UserID)
		assert.Equal(t, "test@example.com", savedUser.Email)
		assert.NotEqual(t, "password123", savedUser.PasswordHash)
	})
}
