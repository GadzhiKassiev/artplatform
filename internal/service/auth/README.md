# Auth Service

Отвечает за регистрацию, аутентификацию и валидацию JWT.

## Возможности

- Регистрация пользователя (`Register`)
- Вход по email/паролю (`Login`)
- Валидация JWT (`ValidateToken`)
- Получение пользователя по ID (`GetUser`)

## Модель
```go
User {
    ID string
    Email string (unique)
    PasswordHash string (bcrypt)
    Role string (STUDENT | TEACHER)
    CreatedAt time.Time
    UpdatedAt time.Time
}
```
## Хранение

- **PostgreSQL**, таблица `users`
- Пароли хешируются через **bcrypt** (cost 10)
- JWT подписывается **HS256**, TTL 24 часа

## gRPC API

См. `proto/auth.proto`.

## Зависимости

- **Postgres** - БД
- **JWT-утилита** - `internal/pkg/jwt`

## Тесты

Планируются:
- Юнит-тесты: `test/`
- Интеграционные: `test/integration/auth/`