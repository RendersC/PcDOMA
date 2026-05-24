# PCDoma — API тесты (Postman / curl)

Запросы для проверки каждого микросервиса. Два способа обращения:

- **Через gateway** (как ходит фронт): `http://localhost:8090/api/v1/...`
- **Напрямую в сервис** (для health-чеков и отладки): порты `9081`–`9084`

| Сервис | Прямой порт | Внешний доступ |
|---|---|---|
| gateway | `8090` | да |
| auth-service | `9081` | да |
| catalog-service | `9082` | да |
| booking-service | `9083` | да |
| payment-service | `9084` | да |
| notification-service | — | только через gateway |

## Переменные окружения Postman

```
base     = http://localhost:8090/api/v1
token    = (заполнится после логина — tokens.access_token)
refresh  = (заполнится после логина — tokens.refresh_token)
pc_id    = (id ПК из GET /catalog/pcs)
booking_id = (id из POST /bookings)
payment_id = (id из POST /payments)
notif_id = (id из GET /notifications)
```

> Совет: в запросе **Login** на вкладке *Tests* добавь авто-сохранение токена:
> ```js
> const j = pm.response.json();
> pm.environment.set("token", j.tokens.access_token);
> pm.environment.set("refresh", j.tokens.refresh_token);
> ```

---

## 1. Auth-service (порт 9081)

### Регистрация
```
POST {{base}}/auth/register
Content-Type: application/json

{
  "name": "Test User",
  "email": "test@pcdoma.kz",
  "password": "secret123",
  "phone": "+77001234567"
}
```

### Логин
```
POST {{base}}/auth/login
Content-Type: application/json

{ "email": "test@pcdoma.kz", "password": "secret123" }
```

### Текущий профиль (защищённый)
```
GET {{base}}/auth/me
Authorization: Bearer {{token}}
```

### Обновить пару токенов
```
POST {{base}}/auth/refresh
Content-Type: application/json

{ "refresh_token": "{{refresh}}" }
```

### Logout (аннулирует refresh-токен в Postgres)
```
POST {{base}}/auth/logout
Content-Type: application/json

{ "refresh_token": "{{refresh}}" }
```

### Health (напрямую)
```
GET http://localhost:9081/health
```

---

## 2. Catalog-service (порт 9082)

### Список ПК (кэшируется в Redis, TTL 300с)
```
GET {{base}}/catalog/pcs
```

### Конкретный ПК
```
GET {{base}}/catalog/pcs/{{pc_id}}
```

### Локации / периферия / сетапы
```
GET {{base}}/catalog/locations
GET {{base}}/catalog/peripherals
GET {{base}}/catalog/setups
```

### Создать ПК (admin, нужен токен — сбрасывает кэш pcs:list:*)
```
POST {{base}}/catalog/pcs
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "name": "RTX 4090 Battlestation",
  "specs": { "cpu": "i9-14900K", "gpu": "RTX 4090", "ram": "64GB" },
  "price_per_hour": 2000,
  "price_per_day": 30000,
  "location": "Astana-Esil",
  "images": []
}
```

### Добавить отзыв (защищённый)
```
POST {{base}}/catalog/pcs/{{pc_id}}/reviews
Authorization: Bearer {{token}}
Content-Type: application/json

{ "rating": 5, "comment": "Топ" }
```

### Список отзывов
```
GET {{base}}/catalog/pcs/{{pc_id}}/reviews
```

### Health (напрямую)
```
GET http://localhost:9082/health
```

---

## 3. Booking-service (порт 9083)

### Свободные слоты (публичный)
```
GET {{base}}/bookings/slots?pc_id={{pc_id}}&date=2026-05-25
```

### Создать бронь (публикует событие в брокер → notification-service)
```
POST {{base}}/bookings
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "booking_type": "pc",
  "pc_id": "{{pc_id}}",
  "location_id": "Astana-Esil",
  "rental_type": "hourly",
  "start_time": "2026-05-25T10:00:00Z",
  "end_time": "2026-05-25T14:00:00Z",
  "peripheral_ids": []
}
```

### Мои брони
```
GET {{base}}/bookings
Authorization: Bearer {{token}}
```

### Одна бронь
```
GET {{base}}/bookings/{{booking_id}}
Authorization: Bearer {{token}}
```

### Отменить бронь (публикует booking.cancelled)
```
DELETE {{base}}/bookings/{{booking_id}}
Authorization: Bearer {{token}}
```

### Очередь работника
```
GET {{base}}/worker/bookings/active
Authorization: Bearer {{token}}

PATCH {{base}}/worker/bookings/{{booking_id}}/accept
Authorization: Bearer {{token}}

PATCH {{base}}/worker/bookings/{{booking_id}}/complete
Authorization: Bearer {{token}}
```

### Админ: все брони + статистика
```
GET {{base}}/admin/bookings
Authorization: Bearer {{token}}

GET {{base}}/admin/bookings/stats
Authorization: Bearer {{token}}
```

### Health (напрямую)
```
GET http://localhost:9083/health
```

---

## 4. Payment-service (порт 9084)

### Создать платёж
```
POST {{base}}/payments
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "booking_id": "{{booking_id}}",
  "amount": 8000,
  "method": "card"
}
```

### Провести оплату
```
POST {{base}}/payments/{{payment_id}}/process
Authorization: Bearer {{token}}
Content-Type: application/json

{ "method": "card" }
```

### Мои платежи
```
GET {{base}}/payments
Authorization: Bearer {{token}}
```

### Один платёж
```
GET {{base}}/payments/{{payment_id}}
Authorization: Bearer {{token}}
```

### Возврат
```
POST {{base}}/payments/{{payment_id}}/refund
Authorization: Bearer {{token}}
```

### Админ: все платежи + статистика
```
GET {{base}}/admin/payments
Authorization: Bearer {{token}}

GET {{base}}/admin/payments/stats
Authorization: Bearer {{token}}
```

### Health (напрямую)
```
GET http://localhost:9084/health
```

---

## 5. Notification-service (только через gateway)

> Уведомления создаёт сам сервис, получая события из брокера (Redis Pub/Sub).
> Чтобы их увидеть: оформи бронь (`POST /bookings`), затем запроси список ниже.

### Мои уведомления
```
GET {{base}}/notifications
Authorization: Bearer {{token}}
```

### Отметить одно прочитанным
```
PATCH {{base}}/notifications/{{notif_id}}/read
Authorization: Bearer {{token}}
```

### Прочитать все
```
PATCH {{base}}/notifications/read-all
Authorization: Bearer {{token}}
```

### Удалить уведомление
```
DELETE {{base}}/notifications/{{notif_id}}
Authorization: Bearer {{token}}
```

---

## Проверка работы брокера сообщений (end-to-end)

1. `POST {{base}}/bookings` — создать бронь.
2. `GET {{base}}/notifications` — появится новое уведомление, созданное notification-service из события `booking.confirmed`.
3. В логах контейнера будет строка `[EMAIL MOCK] To user ...`:
   ```
   docker logs pcdoma-notification --tail 20
   ```

## Проверка кэша каталога (Redis)

1. `GET {{base}}/catalog/pcs` — первый раз читает из MongoDB и кладёт в Redis.
2. Повторный `GET {{base}}/catalog/pcs` — отдаётся из кэша (быстрее).
3. `POST {{base}}/catalog/pcs` (admin) — сбрасывает кэш `pcs:list:*`, следующий список снова идёт в Mongo.

Посмотреть ключи в Redis:
```
docker exec -it pcdoma-redis redis-cli -a secret
> SELECT 1
> KEYS pcs:list:*
```
