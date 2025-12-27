-- +goose Up
CREATE TABLE chirps (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    body TEXT NOT NULL
);

-- +goose Down
DROP TABLE chirps;
