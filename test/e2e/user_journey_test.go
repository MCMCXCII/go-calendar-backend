//go:build e2e

package e2e_test

import (
	"fmt"
	"time"

	"project/pkg/httpclient/authclient"
	"project/pkg/httpclient/eventclient"
)

// Test_NewUserCompleteJourney проходит полный путь нового пользователя:
// регистрация → вход → планирование недели событиями → просмотр
// расписания разными способами → перенос встречи → отмена одной из
// них → выход → попытка воспользоваться системой после выхода.
// Это ровно тот путь, который проходит реальный человек, впервые
// открывший приложение.
func (s *Suite) Test_NewUserCompleteJourney() {
	email := s.uniqueEmail()

	auth := s.newAuthClient()
	events := s.newEventsClient()

	// 1. Регистрация нового пользователя
	userID, err := auth.Register(ctx, email, testPassword)
	s.NoError(err)
	s.NotEmpty(userID)

	// 2. Вход
	token, err := auth.Login(ctx, email, testPassword)
	s.NoError(err)
	s.NotEmpty(token)
	events.SetToken(token)

	// 3. Пользователь планирует рабочую неделю — несколько разных событий
	monday := nextMonday()

	standupID, err := events.CreateEvent(ctx, eventclient.EventInput{
		Title:     "Daily standup",
		Type:      "meeting",
		StartTime: monday.Add(9 * time.Hour),
		EndTime:   monday.Add(9*time.Hour + 15*time.Minute),
	})
	s.NoError(err)

	reviewID, err := events.CreateEvent(ctx, eventclient.EventInput{
		Title:     "Code review",
		Type:      "task",
		StartTime: monday.Add(14 * time.Hour),
		EndTime:   monday.Add(15 * time.Hour),
	})
	s.NoError(err)

	dentistID, err := events.CreateEvent(ctx, eventclient.EventInput{
		Title:      "Стоматолог",
		Type:       "other",
		CustomType: "personal",
		StartTime:  monday.AddDate(0, 0, 2).Add(18 * time.Hour),
		EndTime:    monday.AddDate(0, 0, 2).Add(19 * time.Hour),
	})
	s.NoError(err)

	// 4. Смотрит расписание на понедельник — должно быть два события
	mondayEvents, err := events.ListEvents(ctx, eventclient.ListFilter{
		Day: monday.Format("2006-01-02"),
	})
	s.NoError(err)
	s.Len(mondayEvents, 2)

	// 5. Смотрит расписание на всю неделю — должно быть уже три
	weekEvents, err := events.ListEvents(ctx, eventclient.ListFilter{Week: isoWeekString(monday)})
	s.NoError(err)
	s.Len(weekEvents, 3)

	// 6. Планы меняются — code review переносится на час позже
	err = events.UpdateEvent(ctx, reviewID, eventclient.EventInput{
		Title:     "Code review",
		Type:      "task",
		StartTime: monday.Add(15 * time.Hour),
		EndTime:   monday.Add(16 * time.Hour),
	})
	s.NoError(err)

	updatedReview, err := events.GetEvent(ctx, reviewID)
	s.NoError(err)
	s.True(updatedReview.StartTime.Equal(monday.Add(15 * time.Hour)))

	// 7. Поход к стоматологу отменяется
	err = events.DeleteEvent(ctx, dentistID)
	s.NoError(err)

	_, err = events.GetEvent(ctx, dentistID)
	s.ErrorIs(err, eventclient.ErrNotFound)

	// 8. Расписание на месяц теперь показывает только два оставшихся события
	monthEvents, err := events.ListEvents(ctx, eventclient.ListFilter{
		Month: monday.Format("2006-01"),
	})
	s.NoError(err)
	s.Len(monthEvents, 2)

	// 9. Пользователь заканчивает работу и выходит
	auth.SetToken(token)
	s.NoError(auth.Logout(ctx))

	// 10. Токен отозван — дальнейшие попытки обратиться к своим
	// событиям должны отклоняться, несмотря на то что сам токен
	// формально ещё не истёк по времени
	_, err = events.ListEvents(ctx, eventclient.ListFilter{Day: monday.Format("2006-01-02")})
	s.ErrorIs(err, eventclient.ErrUnauthorized)

	err = auth.Logout(ctx)
	s.ErrorIs(err, authclient.ErrUnauthorized)

	_, err = events.GetEvent(ctx, standupID)
	s.ErrorIs(err, eventclient.ErrUnauthorized)
}

// Test_TwoUsersIndependentJourneys — два разных пользователя одновременно
// пользуются системой, каждый ведёт свой календарь, и ни один не видит
// и не может тронуть события другого ни на одном из шагов.
func (s *Suite) Test_TwoUsersIndependentJourneys() {
	alice := s.registerAndLoginNewUser()
	bob := s.registerAndLoginNewUser()

	day := nextMonday()

	aliceEventID, err := alice.CreateEvent(ctx, eventclient.EventInput{
		Title:     "Alice's meeting",
		Type:      "meeting",
		StartTime: day.Add(10 * time.Hour),
		EndTime:   day.Add(11 * time.Hour),
	})
	s.NoError(err)

	bobEventID, err := bob.CreateEvent(ctx, eventclient.EventInput{
		Title:     "Bob's task",
		Type:      "task",
		StartTime: day.Add(10 * time.Hour),
		EndTime:   day.Add(11 * time.Hour),
	})
	s.NoError(err)

	// у каждого в списке — только своё собственное событие
	aliceEvents, err := alice.ListEvents(ctx, eventclient.ListFilter{Day: day.Format("2006-01-02")})
	s.NoError(err)
	s.Len(aliceEvents, 1)
	s.Equal("Alice's meeting", aliceEvents[0].Title)

	bobEvents, err := bob.ListEvents(ctx, eventclient.ListFilter{Day: day.Format("2006-01-02")})
	s.NoError(err)
	s.Len(bobEvents, 1)
	s.Equal("Bob's task", bobEvents[0].Title)

	// Bob не может ни увидеть, ни удалить событие Alice
	_, err = bob.GetEvent(ctx, aliceEventID)
	s.ErrorIs(err, eventclient.ErrNotFound)

	err = bob.DeleteEvent(ctx, aliceEventID)
	s.ErrorIs(err, eventclient.ErrNotFound)

	// событие Alice всё ещё на месте
	_, err = alice.GetEvent(ctx, aliceEventID)
	s.NoError(err)

	// то же самое верно и в обратную сторону
	_, err = alice.GetEvent(ctx, bobEventID)
	s.ErrorIs(err, eventclient.ErrNotFound)
}

// registerAndLoginNewUser — вспомогательный шаг для многопользовательских
// сценариев: полный цикл "регистрация + вход" одного нового пользователя.
func (s *Suite) registerAndLoginNewUser() *eventclient.Client {
	auth := s.newAuthClient()
	email := s.uniqueEmail()

	_, err := auth.Register(ctx, email, testPassword)
	s.NoError(err)

	token, err := auth.Login(ctx, email, testPassword)
	s.NoError(err)

	events := s.newEventsClient()
	events.SetToken(token)
	return events
}

// nextMonday возвращает полночь ближайшего будущего понедельника —
// удобная фиксированная точка отсчёта для сценариев с целой неделей.
func nextMonday() time.Time {
	now := time.Now().UTC().Truncate(24 * time.Hour)
	offset := (8 - int(now.Weekday())) % 7
	if offset == 0 {
		offset = 7
	}
	return now.AddDate(0, 0, offset)
}

// isoWeekString возвращает ISO-неделю в формате "YYYY-Www", как ожидает
// events-service (см. parseISOWeek в internal/events/service).
func isoWeekString(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%04d-W%02d", year, week)
}
