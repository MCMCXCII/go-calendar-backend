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

func TestService_CreateEvent_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	start := time.Now()
	end := start.Add(time.Hour)

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	var savedEvent domain.Event
	store.EXPECT().
		CreateEvent(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, e domain.Event) error {
			savedEvent = e
			return nil
		})

	svc := service.New(service.Params{Store: store})

	result, err := svc.CreateEvent(ctx, service.CreateEventParams{
		UserID:    userID,
		Title:     "Совещание",
		Type:      domain.EventTypeMeeting,
		StartTime: start,
		EndTime:   end,
	})

	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, result.EventID)
	require.Equal(t, result.EventID, savedEvent.ID)
	require.Equal(t, userID, savedEvent.UserID)
	require.Equal(t, "Совещание", savedEvent.Title)
}

func TestService_CreateEvent_EmptyTitle(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	svc := service.New(service.Params{Store: store})

	_, err := svc.CreateEvent(ctx, service.CreateEventParams{
		Title:     "",
		Type:      domain.EventTypeMeeting,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	})

	require.ErrorIs(t, err, domain.ErrTitleRequired)
}

func TestSerice_CreateEvent_InvalidType(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	svc := service.New(service.Params{Store: store})

	_, err := svc.CreateEvent(ctx, service.CreateEventParams{
		Title:     "Событие",
		Type:      domain.EventType("unknown"),
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	})

	require.ErrorIs(t, err, domain.ErrInvalidType)
}

func TestService_CreateEvent_OtherTypeWithoutCustomType(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	svc := service.New(service.Params{Store: store})

	_, err := svc.CreateEvent(ctx, service.CreateEventParams{
		Title:      "Событие",
		Type:       domain.EventTypeOther,
		CustomType: "", // пусто, хотя type == other
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(time.Hour),
	})

	require.ErrorIs(t, err, domain.ErrCustomTypeRequired)
}

func TestService_CreateEvent_NonOtherTypeWithCustomType(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	svc := service.New(service.Params{Store: store})

	_, err := svc.CreateEvent(ctx, service.CreateEventParams{
		Title:      "Событие",
		Type:       domain.EventTypeMeeting,
		CustomType: "birthday", // лишнее поле для не-other типа
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(time.Hour),
	})

	require.ErrorIs(t, err, domain.ErrCustomTypeNotAllowed)
}

func TestService_CreateEvent_InvalidTimeRange(t *testing.T) {
	ctx := context.Background()
	start := time.Now()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	svc := service.New(service.Params{Store: store})

	_, err := svc.CreateEvent(ctx, service.CreateEventParams{
		Title:     "Событие",
		Type:      domain.EventTypeMeeting,
		StartTime: start,
		EndTime:   start.Add(-time.Hour), // конец раньше начала
	})

	require.ErrorIs(t, err, domain.ErrInvalidTimeRange)
}

func TestService_CreateEvent_StoreErrorIsWrapped(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	dbErr := errors.New("connection refused")
	store.EXPECT().
		CreateEvent(ctx, gomock.Any()).
		Return(dbErr)

	svc := service.New(service.Params{Store: store})

	_, err := svc.CreateEvent(ctx, service.CreateEventParams{
		Title:     "Событие",
		Type:      domain.EventTypeMeeting,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	})

	require.ErrorIs(t, err, dbErr)
}
