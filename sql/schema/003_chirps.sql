-- +goose Up 
ALTER TABLE chirps ADD COLUMN user_id UUID NOT NULL REFERENCES users (id)
ON DELETE CASCADE;

-- +goose Down
ALTER TABLE chirps DROP COLUMN user_id;
