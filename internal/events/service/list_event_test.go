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

func captureQuery(ctx context.Context, t *testing.T, store *mocks.Mockstore) *domain.ListEventsParams {
	t.Helper()

	var captured domain.ListEventsParams
	store.EXPECT().
		ListEvents(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, q domain.ListEventsParams) ([]domain.Event, error) {
			captured = q
			return []domain.Event{}, nil
		})

	return &captured
}

func TestService_ListEvents_ByDay_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)
	captured := captureQuery(ctx, t, store)

	svc := service.New(service.Params{Store: store})

	_, err := svc.ListEvents(ctx, service.ListEventsParams{
		UserID: userID,
		Day:    "2025-04-15",
	})

	require.NoError(t, err)
	require.Equal(t, userID, captured.UserID)
	require.True(t, captured.From.Equal(time.Date(2025, 4, 15, 0, 0, 0, 0, time.UTC)))
	require.True(t, captured.To.Equal(time.Date(2025, 4, 16, 0, 0, 0, 0, time.UTC)))
	require.Equal(t, "day:2025-04-15", captured.CacheKey)
}

func TestService_ListEvents_ByDay_InvalidFormat(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl) // ListEvents не должен вызваться

	svc := service.New(service.Params{Store: store})

	_, err := svc.ListEvents(ctx, service.ListEventsParams{Day: "15-04-2025"})

	require.ErrorIs(t, err, domain.ErrInvalidPeriod)
}

func TestService_ListEvents_ByWeek_Success(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)
	captured := captureQuery(ctx, t, store)

	svc := service.New(service.Params{Store: store})

	_, err := svc.ListEvents(ctx, service.ListEventsParams{Week: "2025-W16"})

	require.NoError(t, err)
	// проверено вручную в разговоре ранее: 2025-W16 -> [2025-04-14, 2025-04-21)
	require.True(t, captured.From.Equal(time.Date(2025, 4, 14, 0, 0, 0, 0, time.UTC)))
	require.True(t, captured.To.Equal(time.Date(2025, 4, 21, 0, 0, 0, 0, time.UTC)))
	require.Equal(t, "week:2025-W16", captured.CacheKey)
}

func TestService_ListEvents_ByWeek_MissingWPrefix(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	svc := service.New(service.Params{Store: store})

	_, err := svc.ListEvents(ctx, service.ListEventsParams{Week: "2025-16"})

	require.ErrorIs(t, err, domain.ErrInvalidPeriod)
}

func TestService_ListEvents_ByWeek_InvalidYear(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	svc := service.New(service.Params{Store: store})

	_, err := svc.ListEvents(ctx, service.ListEventsParams{Week: "abcd-W16"})

	require.ErrorIs(t, err, domain.ErrInvalidPeriod)
}

func TestService_ListEvents_ByWeek_NumberOutOfRange(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	svc := service.New(service.Params{Store: store})

	tests := []string{"2025-W00", "2025-W54"}
	for _, week := range tests {
		_, err := svc.ListEvents(ctx, service.ListEventsParams{Week: week})
		require.ErrorIs(t, err, domain.ErrInvalidPeriod, "week=%s", week)
	}
}

func TestService_ListEvents_ByMonth_Success(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)
	captured := captureQuery(ctx, t, store)

	svc := service.New(service.Params{Store: store})

	_, err := svc.ListEvents(ctx, service.ListEventsParams{Month: "2025-04"})

	require.NoError(t, err)
	require.True(t, captured.From.Equal(time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)))
	require.True(t, captured.To.Equal(time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)))
	require.Equal(t, "month:2025-04", captured.CacheKey)
}

func TestService_ListEvents_ByMonth_InvalidFormat(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	svc := service.New(service.Params{Store: store})

	_, err := svc.ListEvents(ctx, service.ListEventsParams{Month: "April 2025"})

	require.ErrorIs(t, err, domain.ErrInvalidPeriod)
}

func TestService_ListEvents_ByFromTo_Success(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)
	captured := captureQuery(ctx, t, store)

	svc := service.New(service.Params{Store: store})

	_, err := svc.ListEvents(ctx, service.ListEventsParams{
		From: "2025-04-01T00:00:00Z",
		To:   "2025-04-10T00:00:00Z",
	})

	require.NoError(t, err)
	require.True(t, captured.From.Equal(time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)))
	require.True(t, captured.To.Equal(time.Date(2025, 4, 10, 0, 0, 0, 0, time.UTC)))
	require.Empty(t, captured.CacheKey) // произвольный интервал — не кэшируется
}

func TestService_ListEvents_ByFromTo_InvalidFrom(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	svc := service.New(service.Params{Store: store})

	_, err := svc.ListEvents(ctx, service.ListEventsParams{
		From: "not-a-date",
		To:   "2025-04-10T00:00:00Z",
	})

	require.ErrorIs(t, err, domain.ErrInvalidPeriod)
}

func TestService_ListEvents_ByFromTo_InvalidTo(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	svc := service.New(service.Params{Store: store})

	_, err := svc.ListEvents(ctx, service.ListEventsParams{
		From: "2025-04-01T00:00:00Z",
		To:   "not-a-date",
	})

	require.ErrorIs(t, err, domain.ErrInvalidPeriod)
}

func TestService_ListEvents_ByFromTo_FromAfterTo(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	svc := service.New(service.Params{Store: store})

	_, err := svc.ListEvents(ctx, service.ListEventsParams{
		From: "2025-04-10T00:00:00Z",
		To:   "2025-04-01T00:00:00Z",
	})

	require.ErrorIs(t, err, domain.ErrInvalidTimeRange)
}

func TestService_ListEvents_FromToTakesPriorityOverDayWeekMonth(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)
	captured := captureQuery(ctx, t, store)

	svc := service.New(service.Params{Store: store})

	// заданы все параметры сразу — from/to должны победить, day/week/month игнорируются
	_, err := svc.ListEvents(ctx, service.ListEventsParams{
		Day:   "2025-01-01",
		Week:  "2025-W01",
		Month: "2025-01",
		From:  "2025-04-01T00:00:00Z",
		To:    "2025-04-10T00:00:00Z",
	})

	require.NoError(t, err)
	require.True(t, captured.From.Equal(time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)))
	require.Empty(t, captured.CacheKey)
}

func TestService_ListEvents_NoParamsProvided(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl) // ListEvents не должен вызваться

	svc := service.New(service.Params{Store: store})

	_, err := svc.ListEvents(ctx, service.ListEventsParams{})

	require.ErrorIs(t, err, domain.ErrPeriodRequired)
}

func TestService_ListEvents_StoreErrorIsWrapped(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	store := mocks.NewMockstore(ctrl)

	dbErr := errors.New("connection refused")
	store.EXPECT().
		ListEvents(ctx, gomock.Any()).
		Return(nil, dbErr)

	svc := service.New(service.Params{Store: store})

	_, err := svc.ListEvents(ctx, service.ListEventsParams{Day: "2025-04-15"})

	require.ErrorIs(t, err, dbErr)
}
