-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
-- +goose Down
DROP TABLE users;
