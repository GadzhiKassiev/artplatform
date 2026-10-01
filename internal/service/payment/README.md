# Payment Service

Сервис обработки платежей. В прототипе — **заглушка**, всегда возвращает `SUCCESS`.

## Роль в архитектуре

**Core** — обрабатывает транзакции. Вызывается **только** сервисом `purchase`
через gRPC. **Клиент не вызывает payment напрямую** — это внутренний сервис.

## Возможности

- **ProcessPayment** — обработать оплату (заглушка)
- **GetTransaction** — получить транзакцию по ID

## Модель
```go
Transaction {
    ID UUID
    UserID UUID
    CourseID UUID
    Amount float64
    Status PENDING | SUCCESS | FAILED
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

## Хранение

- **PostgreSQL**, таблица `transactions`

## Поток данных

1. purchase.Create() → Payment.ProcessPayment(userID, courseID, amount)

2. Payment → PostgreSQL: INSERT transactions (status=SUCCESS)

3. Payment → Kafka: publish "payment.succeeded" {transactionID, userID, courseID, amount}

4. Payment → purchase: возвращает Transaction

--- в проде ---

6. Внешний провайдер (Stripe, PayPal) обрабатывает реальную оплату

7. Webhook от провайдера → обновление статуса транзакции

## Что в прототипе упрощено

- **Оплата всегда успешна** — нет реального провайдера
- **Нет идемпотентности** — повторный вызов создаст новую транзакцию
- **Нет webhook'ов** — статус устанавливается сразу
- **Нет возвратов** — refund не реализован

## Что в проде

- Интеграция с платёжным провайдером (Stripe, PayPal, ЮKassa)
- PCI DSS compliance
- Idempotency keys
- Webhook handler для обновления статусов
- Saga с компенсациями (refund при неудаче purchase)
- Аудит всех транзакций

## gRPC API

См. `proto/payment.proto`.

## Зависимости

- **Postgres** — хранение транзакций
- **Kafka** — публикация `payment.succeeded`

## Валидация

Protovalidate в `proto/payment.proto`:
- `user_id`, `course_id` — не пустые
- `amount` — больше 0

## Тесты

Планируются:
- Юнит-тесты: `internal/service/payment/test/`
- Интеграционные: `test/integration/payment/`