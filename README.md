# artplatform backend

Бэкенд платформы для продажи обучающего контента. Преподаватели создают курсы, студенты их покупают.

## Стек

- **Go 1.26** - язык
- **gRPC + Protobuf** - общение между сервисами
- **GraphQL (gqlgen)** - API для клиента
- **PostgreSQL 16** - БД
- **Docker + Docker Compose** - запуск

## Архитектура

Client → Gateway (GraphQL) → [gRPC] → auth | course
↓
PostgreSQL

**Что реализовано в прототипе:**
- `auth` - регистрация, логин, валидация JWT
- `course` - CRUD курсов
- `gateway` - GraphQL-фасад над сервисами

**Что не реализовано (следующий этап):**
- `purchase`, `payment`, `media`, `activity` - запланированы, папки созданы
- Kafka - события между сервисами
- MinIO - хранение медиа
- Redis - кэш и rate limit

## Запуск

```bash
docker compose up --build
```

После сборки:

GraphQL Playground: http://localhost:8000

Postgres: localhost:5432 (user: artplatform, pass: artplatform)

auth gRPC: localhost:8001

course gRPC: localhost:8002