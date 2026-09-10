//go:build integration

package integration_test

import (
	"time"

	"github.com/google/uuid"

	"project/pkg/httpclient/eventclient"
)

func (s *Suite) Test_Events_CreateGetUpdateDelete() {
	client := s.registerAndLogin()

	start := time.Now().Truncate(time.Second)
	end := start.Add(time.Hour)

	eventID, err := client.CreateEvent(ctx, eventclient.EventInput{
		Title:     "Совещание",
		Type:      "meeting",
		StartTime: start,
		EndTime:   end,
	})
	s.NoError(err)
	s.NotEqual(uuid.Nil, eventID)

	event, err := client.GetEvent(ctx, eventID)
	s.NoError(err)
	s.Equal("Совещание", event.Title)
	s.Equal("meeting", event.Type)

	newStart := start.Add(2 * time.Hour)
	err = client.UpdateEvent(ctx, eventID, eventclient.EventInput{
		Title:     "Совещание (перенесено)",
		Type:      "meeting",
		StartTime: newStart,
		EndTime:   newStart.Add(time.Hour),
	})
	s.NoError(err)

	updated, err := client.GetEvent(ctx, eventID)
	s.NoError(err)
	s.Equal("Совещание (перенесено)", updated.Title)

	s.NoError(client.DeleteEvent(ctx, eventID))

	_, err = client.GetEvent(ctx, eventID)
	s.ErrorIs(err, eventclient.ErrNotFound)
}

func (s *Suite) Test_Events_CreateWithOtherTypeRequiresCustomType() {
	client := s.registerAndLogin()

	_, err := client.CreateEvent(ctx, eventclient.EventInput{
		Title:     "Событие",
		Type:      "other",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	})

	s.ErrorIs(err, eventclient.ErrValidation)
}

func (s *Suite) Test_Events_ListByDay() {
	client := s.registerAndLogin()
	day := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)

	_, err := client.CreateEvent(ctx, eventclient.EventInput{
		Title:     "Завтрашнее событие",
		Type:      "meeting",
		StartTime: day.Add(10 * time.Hour),
		EndTime:   day.Add(11 * time.Hour),
	})
	s.NoError(err)

	events, err := client.ListEvents(ctx, eventclient.ListFilter{
		Day: day.Format("2006-01-02"),
	})
	s.NoError(err)
	s.Len(events, 1)
	s.Equal("Завтрашнее событие", events[0].Title)
}

func (s *Suite) Test_Events_ListRequiresAtLeastOneFilter() {
	client := s.registerAndLogin()

	_, err := client.ListEvents(ctx, eventclient.ListFilter{})
	s.ErrorIs(err, eventclient.ErrValidation)
}

func (s *Suite) Test_Events_UnauthorizedWithoutToken() {
	client := eventclient.New(eventclient.Config{Host: "localhost", Port: "8080"}) // токен не установлен

	_, err := client.ListEvents(ctx, eventclient.ListFilter{Day: "2026-04-15"})
	s.ErrorIs(err, eventclient.ErrUnauthorized)
}

func (s *Suite) Test_Events_CrossUserIsolation() {
	owner := s.registerAndLogin()
	other := s.registerAndLogin()

	eventID, err := owner.CreateEvent(ctx, eventclient.EventInput{
		Title:     "Приватное событие",
		Type:      "meeting",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	})
	s.NoError(err)

	// чужой пользователь не должен получить доступ ни на чтение, ни на
	// изменение, ни на удаление — везде "not found", а не "forbidden"
	_, err = other.GetEvent(ctx, eventID)
	s.ErrorIs(err, eventclient.ErrNotFound)

	err = other.UpdateEvent(ctx, eventID, eventclient.EventInput{
		Title:     "Попытка перехвата",
		Type:      "meeting",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	})
	s.ErrorIs(err, eventclient.ErrNotFound)

	err = other.DeleteEvent(ctx, eventID)
	s.ErrorIs(err, eventclient.ErrNotFound)

	event, err := owner.GetEvent(ctx, eventID)
	s.NoError(err)
	s.Equal("Приватное событие", event.Title)
}

func (s *Suite) Test_Events_InvalidTimeRange() {
	client := s.registerAndLogin()
	start := time.Now()

	_, err := client.CreateEvent(ctx, eventclient.EventInput{
		Title:     "Событие",
		Type:      "meeting",
		StartTime: start,
		EndTime:   start.Add(-time.Hour),
	})

	s.ErrorIs(err, eventclient.ErrValidation)
}
