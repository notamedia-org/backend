CREATE TABLE user_tokens (
    id TEXT PRIMARY KEY,
    user_id TEXT,
    token TEXT NOT NULL
);
