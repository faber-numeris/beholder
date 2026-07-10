-- +goose Up
CREATE TABLE users
(
    id            CHAR(26) PRIMARY KEY DEFAULT generate_ulid(),
    email         TEXT UNIQUE NOT NULL,
    password_hash TEXT       NOT NULL,
    first_name    TEXT        NOT NULL DEFAULT '',
    last_name     TEXT        NOT NULL DEFAULT '',
    locale        TEXT        NOT NULL DEFAULT '',
    timezone      TEXT        NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ           DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ           DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMPTZ
);

CREATE TABLE user_confirmations
(
    id           CHAR(26) PRIMARY KEY DEFAULT generate_ulid(),
    user_id      CHAR(26) UNIQUE NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token        TEXT            NOT NULL,
    expires_at   TIMESTAMPTZ      NOT NULL,
    confirmed_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ           DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMPTZ           DEFAULT CURRENT_TIMESTAMP
);


-- +goose Down
SELECT 'down SQL query';
drop table users;
drop table user_confirmations;

