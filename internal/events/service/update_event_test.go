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

func TestService_UpdateEvent_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	eventID := uuid.New()
	start := time.Now()
	end := start.Add(time.Hour)

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	var savedEvent domain.Event
	store.EXPECT().
		UpdateEvent(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, e domain.Event) error {
			savedEvent = e
			return nil
		})

	svc := service.New(service.Params{Store: store})

	err := svc.UpdateEvent(ctx, service.UpdateEventParams{
		EventID:   eventID,
		UserID:    userID,
		Title:     "Совещание (перенесено)",
		Type:      domain.EventTypeMeeting,
		StartTime: start,
		EndTime:   end,
	})

	require.NoError(t, err)
	require.Equal(t, eventID, savedEvent.ID)
	require.Equal(t, userID, savedEvent.UserID)
	require.Equal(t, "Совещание (перенесено)", savedEvent.Title)
}

func TestService_UpdateEvent_EmptyTitle(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl) // UpdateEvent не должен вызваться

	svc := service.New(service.Params{Store: store})

	err := svc.UpdateEvent(ctx, service.UpdateEventParams{
		EventID:   uuid.New(),
		UserID:    uuid.New(),
		Title:     "",
		Type:      domain.EventTypeMeeting,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	})

	require.ErrorIs(t, err, domain.ErrTitleRequired)
}

func TestService_UpdateEvent_InvalidType(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	svc := service.New(service.Params{Store: store})

	err := svc.UpdateEvent(ctx, service.UpdateEventParams{
		EventID:   uuid.New(),
		UserID:    uuid.New(),
		Title:     "Событие",
		Type:      domain.EventType("unknown"),
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	})

	require.ErrorIs(t, err, domain.ErrInvalidType)
}

func TestService_UpdateEvent_OtherTypeWithoutCustomType(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	svc := service.New(service.Params{Store: store})

	err := svc.UpdateEvent(ctx, service.UpdateEventParams{
		EventID:    uuid.New(),
		UserID:     uuid.New(),
		Title:      "Событие",
		Type:       domain.EventTypeOther,
		CustomType: "",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(time.Hour),
	})

	require.ErrorIs(t, err, domain.ErrCustomTypeRequired)
}

func TestService_UpdateEvent_NonOtherTypeWithCustomType(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	svc := service.New(service.Params{Store: store})

	err := svc.UpdateEvent(ctx, service.UpdateEventParams{
		EventID:    uuid.New(),
		UserID:     uuid.New(),
		Title:      "Событие",
		Type:       domain.EventTypeMeeting,
		CustomType: "birthday",
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(time.Hour),
	})

	require.ErrorIs(t, err, domain.ErrCustomTypeNotAllowed)
}

func TestService_UpdateEvent_InvalidTimeRange(t *testing.T) {
	ctx := context.Background()
	start := time.Now()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	svc := service.New(service.Params{Store: store})

	err := svc.UpdateEvent(ctx, service.UpdateEventParams{
		EventID:   uuid.New(),
		UserID:    uuid.New(),
		Title:     "Событие",
		Type:      domain.EventTypeMeeting,
		StartTime: start,
		EndTime:   start.Add(-time.Hour),
	})

	require.ErrorIs(t, err, domain.ErrInvalidTimeRange)
}

func TestService_UpdateEvent_NotFound(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	store.EXPECT().
		UpdateEvent(ctx, gomock.Any()).
		Return(domain.ErrEventNotFound)

	svc := service.New(service.Params{Store: store})

	err := svc.UpdateEvent(ctx, service.UpdateEventParams{
		EventID:   uuid.New(),
		UserID:    uuid.New(),
		Title:     "Событие",
		Type:      domain.EventTypeMeeting,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	})

	require.ErrorIs(t, err, domain.ErrEventNotFound)
}

func TestService_UpdateEvent_StoreErrorIsWrapped(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	dbErr := errors.New("connection refused")
	store.EXPECT().
		UpdateEvent(ctx, gomock.Any()).
		Return(dbErr)

	svc := service.New(service.Params{Store: store})

	err := svc.UpdateEvent(ctx, service.UpdateEventParams{
		EventID:   uuid.New(),
		UserID:    uuid.New(),
		Title:     "Событие",
		Type:      domain.EventTypeMeeting,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	})

	require.ErrorIs(t, err, dbErr)
}
