CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name text NOT NULL,
  email text UNIQUE NOT NULL,
  password_hash BLOB NOT NULL,
  activated boolean NOT NULL,
  version integer NOT NULL DEFAULT 1
);
