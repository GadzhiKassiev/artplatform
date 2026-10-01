CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY,
    email         TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role          TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS courses (
    id          UUID PRIMARY KEY,
    author_id   UUID NOT NULL,
    title       TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    price       NUMERIC NOT NULL DEFAULT 0,
    status      TEXT NOT NULL DEFAULT 'DRAFT',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS media_files (
    id           UUID PRIMARY KEY,
    owner_id     UUID NOT NULL,
    course_id    UUID,
    file_name    TEXT NOT NULL,
    content_type TEXT NOT NULL,
    size         BIGINT NOT NULL,
    original_key TEXT NOT NULL,
    preview_key  TEXT,
    status       TEXT NOT NULL DEFAULT 'PENDING',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS transactions (
    id         UUID PRIMARY KEY,
    user_id    UUID NOT NULL,
    course_id  UUID NOT NULL,
    amount     NUMERIC NOT NULL,
    status     TEXT NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS purchases (
    id             UUID PRIMARY KEY,
    user_id        UUID NOT NULL,
    course_id      UUID NOT NULL,
    amount         NUMERIC NOT NULL,
    status         TEXT NOT NULL DEFAULT 'PENDING',
    transaction_id UUID,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS course_stats (
    course_id      UUID PRIMARY KEY,
    view_count     BIGINT NOT NULL DEFAULT 0,
    purchase_count BIGINT NOT NULL DEFAULT 0,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_purchases_user ON purchases(user_id);
CREATE INDEX IF NOT EXISTS idx_purchases_course ON purchases(course_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_purchases_user_course ON purchases(user_id, course_id);

CREATE INDEX IF NOT EXISTS idx_transactions_user ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_course ON transactions(course_id);

CREATE INDEX IF NOT EXISTS idx_media_owner ON media_files(owner_id);
CREATE INDEX IF NOT EXISTS idx_media_course ON media_files(course_id);
CREATE INDEX IF NOT EXISTS idx_media_status ON media_files(status);

CREATE INDEX IF NOT EXISTS idx_courses_author ON courses(author_id);