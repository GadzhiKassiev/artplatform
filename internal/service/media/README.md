# Media Service

Сервис для загрузки и хранения медиа-файлов (обложки курсов, изображения).

## Роль в архитектуре

**Adapter** — оборачивает MinIO (S3-совместимое хранилище) и предоставляет
другим сервисам чистый интерфейс для работы с файлами.

Сервис `course` не знает про MinIO, S3, presigned URLs. Он просто говорит
`media.UploadMedia(...)`. Все детали работы с хранилищем — внутри media.

## Возможности

- **UploadMedia** — загрузка файла в MinIO
- **GetMedia** — получение метаданных файла по ID
- **GetPresignedURL** — генерация временной ссылки для скачивания

## Модель
```go
MediaFile {
    ID UUID
    OwnerID UUID (кто загрузил)
    CourseID UUID (опционально — к какому курсу относится)
    FileName string
    ContentType string (image/jpeg, image/png, ...)
    Size int64
    OriginalKey string (путь в MinIO: media/{id}/original/{filename})
    PreviewKey string (путь превью: media/{id}/preview/{filename})
    Status PENDING | PROCESSING | READY | FAILED
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

## Статусы

- **PENDING** — запись создана, файл ещё не загружен (не используется в текущем флоу)
- **PROCESSING** — файл загружен в MinIO, ждёт генерации превью
- **READY** — превью готово, файл полностью доступен
- **FAILED** — ошибка при обработке (например, невалидный файл)

## Хранение

- **PostgreSQL** — метаданные файлов (`media_files`)
- **MinIO** — сами файлы
  - Оригинал: `media/{id}/original/{filename}`
  - Превью: `media/{id}/preview/{filename}`

## Поток данных

1. Gateway → Media.UploadMedia(ownerID, fileName, contentType, content)

2. Media → MinIO: PutObject (сохраняет оригинал)

3. Media → PostgreSQL: INSERT media_files (status=PROCESSING)

4. Media → Kafka: publish "media.uploaded" {mediaID, originalKey, contentType}

5. Media → Gateway: возвращает MediaFile

--- асинхронно ---

6. media-worker → Kafka: читает "media.uploaded"

7. media-worker → MinIO: GetObject (скачивает оригинал)

8. media-worker: генерирует превью (ресайз)

9. media-worker → MinIO: PutObject (сохраняет превью)

10. media-worker → PostgreSQL: UPDATE status=READY, preview_key=...


## Presigned URLs

- **TTL 1 час** — время жизни ссылки
- Если превью готово (status=READY) — отдаём **превью**
- Если нет — отдаём **оригинал**

## Права доступа

- **Загружать** — любой аутентифицированный пользователь
- **Скачивать presigned URL** — только владелец файла (в прототипе)
- В проде: доступ через enrollment (куплен ли курс), роль админа

## gRPC API

См. `proto/media.proto`.

## Зависимости

- **Postgres** — метаданные
- **MinIO** — хранилище файлов
- **Kafka** — публикация событий

## Валидация

Через Protovalidate в `proto/media.proto`:
- `content_type` — только image-типы (`image/jpeg`, `image/png`, `image/webp`, `image/gif`)
- `content` — минимум 1 байт
- `file_name`, `owner_id` — не пустые

## Тесты

Планируются:
- Юнит-тесты: `internal/service/media/test/`
- Интеграционные: `test/integration/media/`

## Вне скоупа

- Видео-файлы (только картинки)
- Транскодирование
- Множественная загрузка одним запросом
- Удаление файлов (есть метод в MinIO-клиенте, не подключён к API)