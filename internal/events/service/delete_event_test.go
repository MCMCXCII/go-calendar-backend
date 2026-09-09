package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"project/internal/events/domain"
	"project/internal/events/service"
	"project/internal/events/service/mocks"
)

func TestService_DeleteEvent_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	eventID := uuid.New()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	store.EXPECT().
		DeleteEvent(ctx, userID, eventID).
		Return(nil)

	svc := service.New(service.Params{Store: store})

	err := svc.DeleteEvent(ctx, userID, eventID)

	require.NoError(t, err)
}

func TestService_DeleteEvent_NotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	eventID := uuid.New()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	store.EXPECT().
		DeleteEvent(ctx, userID, eventID).
		Return(domain.ErrEventNotFound)

	svc := service.New(service.Params{Store: store})

	err := svc.DeleteEvent(ctx, userID, eventID)

	require.ErrorIs(t, err, domain.ErrEventNotFound)
}

func TestService_DeleteEvent_StoreErrorIsWrapped(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	eventID := uuid.New()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	dbErr := errors.New("connection refused")
	store.EXPECT().
		DeleteEvent(ctx, userID, eventID).
		Return(dbErr)

	svc := service.New(service.Params{Store: store})

	err := svc.DeleteEvent(ctx, userID, eventID)

	require.ErrorIs(t, err, dbErr)
	require.NotErrorIs(t, err, domain.ErrEventNotFound)
}
