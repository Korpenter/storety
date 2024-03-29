-- +goose Up
CREATE TABLE IF NOT EXISTS users (
     id uuid UNIQUE NOT NULL PRIMARY KEY,
     username text UNIQUE NOT NULL,
     password text NOT NULL,
     salt     text NOT NULL,
     created_at timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS data (
    id uuid UNIQUE NOT NULL PRIMARY KEY,
    user_id uuid NOT NULL,
    name text,
    type varchar(10) CHECK (type IN ('Card', 'Cred', 'Binary', 'Text')),
    content bytea,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
    deleted boolean NOT NULL DEFAULT false,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sessions (
    id uuid UNIQUE NOT NULL PRIMARY KEY,
    user_id uuid NOT NULL,
    auth_token text NOT NULL,
    refresh_token text NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "data";
DROP TABLE IF EXISTS "sessions";
-- +goose StatementEnd