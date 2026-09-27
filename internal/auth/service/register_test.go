package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"project/internal/auth/domain"
	"project/internal/auth/service"
	"project/internal/auth/service/mocks"
)

func TestService_Register_Success(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockStore(ctrl)

	var savedUser domain.User
	store.EXPECT().
		CreateUser(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, user domain.User) error {
			savedUser = user
			return nil
		})

	svc := service.New(service.Params{Store: store})

	result, err := svc.Register(ctx, service.RegisterParams{
		Email:    "Test@example.com",
		Password: "password123",
	})

	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, result.UserID)
	require.Equal(t, result.UserID, savedUser.ID)
	require.Equal(t, "test@example.com", savedUser.Email)
	require.NotEmpty(t, savedUser.PasswordHash)
	require.NotEqual(t, "password123", savedUser.PasswordHash)
}

func TestService_Register_EmailAlreadyExists(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockStore(ctrl)

	store.EXPECT().
		CreateUser(ctx, gomock.Any()).
		Return(domain.ErrEmailAlreadyExists)

	svc := service.New(service.Params{Store: store})

	_, err := svc.Register(ctx, service.RegisterParams{Email: "test@example.com", Password: "password123"})

	require.ErrorIs(t, err, domain.ErrEmailAlreadyExists)
}
