package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"project/internal/auth/domain"
	"project/internal/auth/service"
	"project/internal/auth/service/mocks"
)

func newTestUser(t *testing.T, password string) domain.User {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	return domain.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: string(hash),
	}
}

func TestService_Login_Success(t *testing.T) {
	ctx := context.Background()
	existingUser := newTestUser(t, "correct-password")

	ctrl := gomock.NewController(t)
	store := mocks.NewMockStore(ctrl)
	tokenBuilder := mocks.NewMockTokenBuilder(ctrl)

	store.EXPECT().
		GetUser(ctx, "test@example.com").
		Return(existingUser, nil)

	tokenBuilder.EXPECT().
		BuildAccessToken(existingUser.ID, 24*time.Hour).
		Return("signed.jwt.token", nil)

	svc := service.New(service.Params{
		Store:           store,
		Token:           tokenBuilder,
		TokenExpiration: 24 * time.Hour,
	})

	result, err := svc.Login(ctx, service.LoginParams{
		Email:    "Test@Example.com",
		Password: "correct-password",
	})

	require.NoError(t, err)
	require.Equal(t, "signed.jwt.token", result.AccessToken)
}

func TestService_Login_EmptyEmailOrPassword(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockStore(ctrl)

	svc := service.New(service.Params{Store: store})

	_, err := svc.Login(ctx, service.LoginParams{Email: "", Password: ""})

	require.ErrorIs(t, err, service.ErrInvalidCredentials)
}

func TestService_Login_UserNotFound(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockStore(ctrl)

	store.EXPECT().
		GetUser(ctx, "unknown@example.com").
		Return(domain.User{}, domain.ErrUserNotFound)

	svc := service.New(service.Params{Store: store})

	_, err := svc.Login(ctx, service.LoginParams{Email: "unknown@example.com", Password: "whatever"})

	require.ErrorIs(t, err, service.ErrInvalidCredentials)
	require.NotErrorIs(t, err, domain.ErrUserNotFound)
}

func TestService_Login_WrongPassword(t *testing.T) {
	ctx := context.Background()
	existingUser := newTestUser(t, "correct-password")

	ctrl := gomock.NewController(t)
	store := mocks.NewMockStore(ctrl)

	store.EXPECT().
		GetUser(ctx, "test@example.com").
		Return(existingUser, nil)

	svc := service.New(service.Params{Store: store})

	_, err := svc.Login(ctx, service.LoginParams{Email: "test@example.com", Password: "wrong-password"})

	require.ErrorIs(t, err, service.ErrInvalidCredentials)
}

func TestService_Login_BuildAccessTokenError(t *testing.T) {
	ctx := context.Background()
	existingUser := newTestUser(t, "correct-password")

	ctrl := gomock.NewController(t)
	store := mocks.NewMockStore(ctrl)
	tokenBuilder := mocks.NewMockTokenBuilder(ctrl)

	store.EXPECT().
		GetUser(ctx, "test@example.com").
		Return(existingUser, nil)

	tokenErr := errors.New("sign token: key error")
	tokenBuilder.EXPECT().
		BuildAccessToken(existingUser.ID, gomock.Any()).
		Return("", tokenErr)

	svc := service.New(service.Params{Store: store, Token: tokenBuilder})

	_, err := svc.Login(ctx, service.LoginParams{Email: "test@example.com", Password: "correct-password"})

	require.ErrorIs(t, err, tokenErr)
}
