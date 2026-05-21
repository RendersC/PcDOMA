# PCDoma — сервис аренды ПК (Астана)

Микросервисный сервис аренды игровых и рабочих ПК с доставкой по Астане.
Финальный проект по дисциплине **Advanced Programming 2**.

Пользователь собирает конфигурацию (ПК + периферия) или берёт готовый сетап,
оплачивает аренду, а работник на точке принимает заказ, выдаёт и завершает его.

---

## Архитектура

```
                         React SPA (nginx) :3000
                                  │
                                  ▼
                        API Gateway :8090  ──── единственный публичный вход
              JWT-валидация · reverse proxy · CORS · проброс X-User-*
        ┌──────────┬──────────┬───────────┬───────────┬──────────────┐
        ▼          ▼          ▼           ▼           ▼              
   auth :8081  catalog :8082  booking :8083  payment :8084  notification :8085
        │          │              │            │              │
   PostgreSQL   MongoDB       PostgreSQL   PostgreSQL       MongoDB
              + Redis(cache)        │
                                    └── Redis Pub/Sub ──► notification
```

- **Sync (REST)**: gateway → auth (`/internal/validate-token`); booking → catalog (`/internal/pcs/:id/availability`, `/status`) и payment (`/internal/payments`).
- **Async (Redis Pub/Sub)**: booking/payment публикуют события в каналы `booking.events`, `worker.events`, `payment.events`; notification-service подписан и пишет уведомления в MongoDB.

### Принципы (Clean Architecture)

Каждый сервис: `handler → service → repository → domain`. Слой `service`
зависит от **интерфейсов** репозиториев и клиентов (см. `internal/service/*.go`),
что даёт изоляцию БД (database-per-service) и юнит-тестируемость через моки.

---

## Стек

| Слой | Технология |
|---|---|
| Backend | Go 1.23 + Gin |
| Транзакционная БД | PostgreSQL 16 (auth, booking, payment) |
| Документная БД | MongoDB 7 (catalog, notifications) |
| Кэш / Pub-Sub | Redis 7 |
| Frontend | React + TypeScript + Vite + Tailwind + TanStack Query + Zustand |
| Контейнеризация | Docker + Docker Compose |
| API-документация | Swagger (swaggo) |
| Тесты | testify / testify/mock |

---

## Быстрый старт

```bash
cp .env.example .env          # при необходимости поправьте секреты
docker compose up --build -d  # поднимет 13 контейнеров (миграции применятся автоматически)
```

Наполнить каталог тестовым контентом (18 ПК: 13 реальных + 5 «пасхалок»,
периферия и сетапы):

```powershell
pwsh ./scripts/seed-catalog.ps1
```

Открыть приложение: **http://localhost:3000**

### Демо-доступ

| Роль | Email | Пароль |
|---|---|---|
| admin | `test@pcdoma.kz` | `password123` |

Обычного пользователя можно зарегистрировать прямо в UI.

---

## Порты

| Сервис | Внутри сети | На хосте |
|---|---|---|
| Frontend (nginx) | 80 | **3000** |
| API Gateway | 8080 | **8090** |
| auth-service | 8081 | 9081 |
| catalog-service | 8082 | 9082 |
| booking-service | 8083 | 9083 |
| payment-service | 8084 | 9084 |
| notification-service | 8085 | — |
| PostgreSQL / MongoDB / Redis | 5432 / 27017 / 6379 | — |

Все клиентские запросы идут через gateway: `http://localhost:8090/api/v1/...`

---

## Основные эндпоинты (через gateway, префикс `/api/v1`)

| Метод | Путь | Назначение |
|---|---|---|
| POST | `/auth/register`, `/auth/login`, `/auth/refresh` | регистрация/логин/refresh |
| GET | `/auth/me` | профиль |
| GET | `/catalog/pcs`, `/catalog/pcs/:id` | каталог ПК (кэш Redis 5 мин) |
| GET | `/catalog/peripherals`, `/catalog/setups` | периферия и готовые сетапы |
| POST/PUT/DELETE | `/catalog/pcs` … | CRUD каталога (admin) |
| POST | `/bookings` | создать бронь (custom/setup) → оплата → очередь работника |
| GET | `/bookings`, `/bookings/:id` | мои брони |
| DELETE | `/bookings/:id` | отмена |
| GET/PATCH | `/worker/bookings…` | очередь и статусы для работника |
| GET | `/admin/bookings`, `/admin/payments` | админ-обзор |
| GET | `/notifications` | уведомления пользователя |

### Жизненный цикл брони

```
pending_payment → paid → pending_worker → accepted → delivering → active → completed
                                   ↓                                  
                                rejected / cancelled
```

---

## Swagger UI

| Сервис | URL |
|---|---|
| auth | http://localhost:9081/swagger/index.html |
| catalog | http://localhost:9082/swagger/index.html |
| booking | http://localhost:9083/swagger/index.html |
| payment | http://localhost:9084/swagger/index.html |

Спецификации генерируются `swag init` на этапе сборки Docker-образа.

---

## Тесты

Юнит-тесты на слое `service` (моки репозиториев/клиентов через `testify/mock`),
плюс тесты middleware gateway и маппинга событий notification.

```bash
# по каждому сервису
cd auth-service        && go test ./...
cd catalog-service     && go test ./internal/service/...
cd booking-service     && go test ./internal/service/...
cd payment-service     && go test ./internal/service/...
cd notification-service && go test ./internal/events/...
cd gateway             && go test ./internal/middleware/...
```

Покрытие по сервисам:

| Сервис | Что проверяется |
|---|---|
| auth | register/login/валидация JWT, неверный пароль, занятый email |
| catalog | создание/листинг/доступность ПК, пересчёт рейтинга по отзыву |
| booking | расчёт цены, недоступный ПК, переходы статусов worker, отмена, доступы |
| payment | дефолт валюты, mock-обработка, повторная обработка, возврат |
| notification | маппинг booking/worker/payment событий в уведомления |
| gateway | auth-middleware: без токена / неверный формат / валидный / невалидный |

---

## Структура

```
PCDoma/
├── docker-compose.yml        # 13 контейнеров
├── .env.example
├── Makefile
├── gateway/                  # reverse proxy + JWT middleware
├── auth-service/             # PostgreSQL, JWT, RBAC
├── catalog-service/          # MongoDB + Redis cache (ПК/периферия/сетапы)
├── booking-service/          # PostgreSQL, оркестрация sync+async
├── payment-service/          # PostgreSQL, mock-оплата
├── notification-service/     # MongoDB, Redis subscriber
├── frontend/                 # React SPA
└── scripts/
    ├── init-databases.sql    # создаёт 4 БД в PostgreSQL
    └── seed-catalog.ps1      # наполнение каталога
```

---

## Полезное

```bash
docker compose ps                    # статус контейнеров
docker compose logs -f booking-service
docker compose down                  # остановить
docker compose down -v               # остановить и удалить тома (сброс данных)
curl http://localhost:8090/health    # health gateway
```
