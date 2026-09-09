package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"project/internal/auth/service"
	"project/internal/auth/service/mocks"
)

func TestService_Logout_Success(t *testing.T) {
	ctx := context.Background()
	expiresAt := time.Now().Add(2 * time.Hour)

	ctrl := gomock.NewController(t)
	blacklist := mocks.NewMockBlacklist(ctrl)

	blacklist.EXPECT().
		Revoke(ctx, "jti-123", gomock.Any()).
		DoAndReturn(func(_ context.Context, _ string, ttl time.Duration) error {
			// TTL должен быть примерно равен оставшемуся времени жизни токена
			require.InDelta(t, 2*time.Hour, ttl, float64(time.Second))
			return nil
		})

	svc := service.New(service.Params{Blacklist: blacklist})

	err := svc.Logout(ctx, service.LogoutParams{
		TokenID:   "jti-123",
		ExpiresAt: expiresAt,
	})

	require.NoError(t, err)
}

func TestService_Logout_EmptyTokenID(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	blacklist := mocks.NewMockBlacklist(ctrl) // Revoke не должен вызваться

	svc := service.New(service.Params{Blacklist: blacklist})

	err := svc.Logout(ctx, service.LogoutParams{
		TokenID:   "",
		ExpiresAt: time.Now().Add(time.Hour),
	})

	require.Error(t, err)
}

func TestService_Logout_AlreadyExpiredToken(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	blacklist := mocks.NewMockBlacklist(ctrl) // Revoke не должен вызваться

	svc := service.New(service.Params{Blacklist: blacklist})

	err := svc.Logout(ctx, service.LogoutParams{
		TokenID:   "jti-123",
		ExpiresAt: time.Now().Add(-time.Hour), // уже в прошлом
	})

	require.ErrorIs(t, err, service.ErrTokenExpired)
}

func TestService_Logout_BlacklistErrorIsWrapped(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	blacklist := mocks.NewMockBlacklist(ctrl)

	revokeErr := errors.New("redis connection failed")
	blacklist.EXPECT().
		Revoke(ctx, "jti-123", gomock.Any()).
		Return(revokeErr)

	svc := service.New(service.Params{Blacklist: blacklist})

	err := svc.Logout(ctx, service.LogoutParams{
		TokenID:   "jti-123",
		ExpiresAt: time.Now().Add(time.Hour),
	})

	require.ErrorIs(t, err, revokeErr)
}
