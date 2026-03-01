CREATE TABLE IF NOT EXISTS movies (
    id         SERIAL PRIMARY KEY,
    genre      TEXT NOT NULL DEFAULT '',
    budget     BIGINT NOT NULL DEFAULT 0,
    title      TEXT NOT NULL,
    hero       TEXT NOT NULL DEFAULT '',
    heroine    TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_movies_title ON movies (title);
