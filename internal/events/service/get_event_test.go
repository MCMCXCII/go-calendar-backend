package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"project/internal/events/domain"
	"project/internal/events/service"
	"project/internal/events/service/mocks"
)

func TestService_GetEvent_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	eventID := uuid.New()

	expected := domain.Event{
		ID:        eventID,
		UserID:    userID,
		Title:     "Совещание",
		Type:      domain.EventTypeMeeting,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	}

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	store.EXPECT().
		GetEvent(ctx, userID, eventID).
		Return(expected, nil)

	svc := service.New(service.Params{Store: store})

	result, err := svc.GetEvent(ctx, userID, eventID)

	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestService_GetEvent_NotFound(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	eventID := uuid.New()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	store.EXPECT().
		GetEvent(ctx, userID, eventID).
		Return(domain.Event{}, domain.ErrEventNotFound)

	svc := service.New(service.Params{Store: store})

	_, err := svc.GetEvent(ctx, userID, eventID)

	require.ErrorIs(t, err, domain.ErrEventNotFound)
}

func TestService_GetEvent_StoreErrorIsWrapped(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	eventID := uuid.New()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	dbErr := errors.New("connection refused")
	store.EXPECT().
		GetEvent(ctx, userID, eventID).
		Return(domain.Event{}, dbErr)

	svc := service.New(service.Params{Store: store})

	_, err := svc.GetEvent(ctx, userID, eventID)

	require.ErrorIs(t, err, dbErr)
	require.NotErrorIs(t, err, domain.ErrEventNotFound)
}
