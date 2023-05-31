CREATE TABLE IF NOT EXISTS movies(
    id INTEGER PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT "2023",
    title TEXT not NULL,
    year INTEGER NOT NULL CHECK (year BETWEEN 1888 AND 2023),
    runtime INTEGER CHECK (runtime >90),
    genres TEXT NOT NULL,
    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS movies_title_idx ON movies(title);
CREATE INDEX IF NOT EXISTS movies_genres_idx ON movies(genres);