# Purchase Service

Сервис покупок курсов. Оркестрирует процесс покупки: проверяет, не куплено ли,
вызывает `payment`, создаёт запись о покупке.

## Роль в архитектуре

**Core** — хранит историю покупок. Вызывается клиентом через gateway.

## Возможности

- **CreatePurchase** — купить курс
- **GetPurchase** — получить покупку по ID
- **GetUserPurchases** — все покупки пользователя
- **CheckAccess** — проверка, куплен ли курс

## Модель
```go
Purchase {
    ID UUID
    UserID UUID
    CourseID UUID
    Amount float64
    Status PENDING | COMPLETED | REFUNDED | CANCELLED
    TransactionID UUID (ссылка на payment)
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

## Хранение

- **PostgreSQL**, таблица `purchases`
- **Уникальный индекс** на `(user_id, course_id)` — нельзя купить курс дважды

## Поток данных

1. Клиент → Gateway → Purchase.CreatePurchase(courseID, amount)

2. Purchase → PostgreSQL: проверка "уже куплено?" (HasAccess)

3. Purchase → Payment.ProcessPayment (gRPC)

4. Payment → Kafka: "payment.succeeded"

5. Purchase → PostgreSQL: INSERT purchases (status=COMPLETED)

6. Purchase → Kafka: "purchase.created"

7. Purchase → Gateway → Клиент: Purchase

## Ключевые моменты

**Cross-service вызов:** purchase зависит от payment через gRPC. В проде — 
с saga и компенсациями.

**Идемпотентность:** уникальный индекс `(user_id, course_id)` предотвращает 
двойные покупки. При попытке — ошибка `ALREADY_PURCHASED`.

**Управление доступом:** `CheckAccess` возвращает, куплен ли курс. 
Используется другими сервисами (например, media для доступа к контенту).

## gRPC API

См. `proto/purchase.proto`.

## Зависимости

- **Postgres** — хранение покупок
- **Kafka** — публикация `purchase.created`
- **Payment** (gRPC) — проведение оплаты

## Валидация

Protovalidate в `proto/purchase.proto`:
- `user_id`, `course_id` — не пустые
- `amount` — больше 0

## Тесты

Планируются:
- Юнит-тесты: `internal/service/purchase/test/`
- Интеграционные: `test/integration/purchase/`

## Что в проде

- **Saga** с компенсацией (refund, если purchase упал после payment)
- **Idempotency keys** для повторных запросов
- **Outbox pattern** для надёжной публикации событий
- **Возвраты** (refund) с полным флоу
- **Coupon/promo** интеграция