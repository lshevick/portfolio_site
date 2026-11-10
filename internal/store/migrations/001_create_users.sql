-- 001_create_users.sql
CREATE TABLE users (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    username      TEXT NOT NULL,
    email         TEXT NOT NULL UNIQUE,   -- unique guard -> ErrDuplicateEmails
    password      TEXT NOT NULL,          -- bcrypt hash only; never plaintext
    created_at    DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);
