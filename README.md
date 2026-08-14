# go-calendar-backend

┌──────────────┐
   Клиент ────────▶ │ auth-service │ ───▶ PostgreSQL (auth_db)
                    │   :8081      │ ───▶ Redis (blacklist токенов)
                    └──────────────┘

                    ┌──────────────┐
   Клиент ────────▶ │ events-service│ ───▶ PostgreSQL (events_db)
   (с JWT)          │   :8080       │ ───▶ Redis (кэш списков событий)
                    └──────────────┘

Быстрый старт:
git clone <url-репозитория>
cd go-calendar-backend/deployments
docker compose up --build
