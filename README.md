# go-calendar-backend

Backend календаря на микросервисной архитектуре. Два независимых сервиса — аутентификация и события — с общей JWT-авторизацией, Redis-кэшированием и изоляцией данных на уровне БД.

## Стек
- Go 1.25 
- chi 
- PostgreSQL 
- Redis 
- golang-migrate 
- JWT 
- bcrypt 
- Docker Compose

## Архитектура
- **auth-service** (`:8081`) — регистрация, логин, logout, своя БД
- **events-service** (`:8080`) — CRUD событий с фильтром по дню/неделе/месяцу, своя БД
- Сервисы не общаются напрямую: JWT валидируется в events-service самостоятельно (подпись + чёрный список в Redis)
- Redis: чёрный список токенов (auth) + кэш списков событий с инвалидацией (events)

## Запуск

```bash
cd deployments
docker compose up --build
```

