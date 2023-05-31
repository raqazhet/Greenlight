CREATE TABLE IF NOT EXISTS tokens (
    hash BLOB PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expiry timestamp NOT NULL,
    scope text NOT NULL
);