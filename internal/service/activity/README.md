# Activity Service

Сервис статистики: считает просмотры и покупки курсов.

## Роль в архитектуре

**Adapter** — оборачивает Redis (счётчики) и Kafka (события).
Не хранит бизнес-данные, только счётчики.

## Возможности

- **GetCourseStats** — получить статистику курса (просмотры, покупки)
- **RecordView** — записать просмотр курса
- **Kafka consumers** — автоматически обновляют счётчики из событий

## Модель
```go
CourseStats {
    CourseID UUID
    ViewCount int64
    PurchaseCount int64
    UpdatedAt time.Time
}
```
## Хранение

- **Redis** — счётчики (`stats:course:{id}:views`, `stats:course:{id}:purchases`)
- **TTL 24 часа** — счётчики сбрасываются, если не обновляются
- В проде: периодический flush в Postgres для durability

## Поток данных

**Запись (события):**

1. Kafka topic "course.viewed" → activity.handleCourseViewed → INCR stats:...:views

2. Kafka topic "purchase.created" → activity.handlePurchaseCreated → INCR stats:...:purchases

**Чтение:**

Клиент → Gateway → Activity.GetCourseStats → Redis GET → stats

## Consumer group

Оба consumer'а в одной группе `activity-worker`. Несколько инстансов activity 
**делят нагрузку** — каждое событие попадёт только одному инстансу.

## Почему Redis, а не Postgres

**Счётчики — высокочастотные.** Каждый просмотр курса = INCR. 
Postgres не выдержит такого количества UPDATE'ов.

Redis даёт:
- **Атомарный INCR** — быстрый и безопасный
- **TTL** — автоматическая очистка
- **Pipeline** — batch операций

## Что в проде

- **Периодический flush в Postgres** — для durability
- **Агрегация по времени** — просмотры за день/неделю/месяц
- **Prometheus метрики** — для мониторинга
- **Redis Cluster** — для масштабирования

## gRPC API

См. `proto/activity.proto`.

## Зависимости

- **Redis** — счётчики
- **Kafka** — чтение событий

## Валидация

Protovalidate в `proto/activity.proto`:
- `course_id` — не пустой

## Тесты

Планируются:
- Юнит-тесты: `internal/service/activity/test/`
- Интеграционные: `test/integration/activity/`