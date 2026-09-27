# Course Service

Отвечает за CRUD курсов.

## Возможности

- Создание курса (`CreateCourse`)
- Получение курса по ID (`GetCourse`)
- Список курсов (`GetCourses`) - всех опубликованных или по автору
- Обновление курса (`UpdateCourse`) - только автор
- Удаление курса (`DeleteCourse`) - только автор

## Модель
```go
Course {
    ID string
    AuthorID string (ссылка на auth.UserID)
    Title string
    Description string
    Price float64
    Status string (DRAFT | PUBLISHED)
    CreatedAt time.Time
    UpdatedAt time.Time
}
```
## Хранение

- **PostgreSQL**, таблица `courses`
- Индекс по `author_id`

## Правила

- Создать курс может только аутентифицированный пользователь
- Изменять/удалять курс может **только автор** (проверка по `RequesterID`)
- `getCourses` без фильтра возвращает только `PUBLISHED`
- `getCourses(authorId)` возвращает все курсы автора (включая `DRAFT`)

## Связь с другими сервисами

- **AuthorID** - это `UserID` из auth-сервиса. Хранится как `string`, **без FK** (разные сервисы, разные БД)
- Валидация «автор существует и имеет роль TEACHER» - на уровне **gateway** (перед вызовом)

## gRPC API

См. `proto/course.proto`.

## Тесты

Планируются:
- Юнит-тесты: `test/`
- Интеграционные: `test/integration/course/`