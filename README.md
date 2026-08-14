# go-calendar-backend

Backend для календаря на Go: два независимых микросервиса — аутентификация и события, с общей авторизацией через JWT.

## Стек
Go 1.25 · chi · PostgreSQL (pgx) · Redis · golang-migrate · JWT · bcrypt · Docker Compose

## Архитектура
- **auth-service** (`:8081`) — регистрация, логин, logout, своя БД
- **events-service** (`:8080`) — CRUD событий с фильтром по дню/неделе/месяцу, своя БД
- Сервисы не общаются напрямую: JWT валидируется в events-service самостоятельно (подпись + чёрный список в Redis)
- Redis: чёрный список токенов (auth) + кэш списков событий с инвалидацией (events)

## Запуск

```bash
cd deployments
cp .env.example .env
docker compose up --build
```

## API

**Auth** — `/api/v1/auth`: `POST /register`, `POST /login`, `POST /logout`, `GET /me`

**Events** — `/api/v1/events` (требует `Authorization: Bearer <token>`):
`POST /`, `GET /?day=YYYY-MM-DD|week=YYYY-Www|month=YYYY-MM`, `GET/{id}`, `PUT /{id}`, `DELETE /{id}`

## Тестирование

Коллекция Postman — в `postman/`.
