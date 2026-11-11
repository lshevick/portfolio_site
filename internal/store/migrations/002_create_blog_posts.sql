-- 002_create_blog_posts.sql
CREATE TABLE blog_posts (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    title       TEXT NOT NULL,
    slug        TEXT NOT NULL UNIQUE,       -- URL-friendly: /blog/<slug>
    author_id   INTEGER NOT NULL REFERENCES users(id),  -- FK to users.id
    author_name TEXT NOT NULL,              -- denormalized for display; avoids a JOIN in templates
    body        TEXT NOT NULL,              -- post content (HTML or markdown)
    published   BOOLEAN NOT NULL DEFAULT 0, -- 0 = draft, 1 = live
    created_at  DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP),
    updated_at  DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);
