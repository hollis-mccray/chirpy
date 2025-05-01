-- +goose Up
ALTER TABLE users
ADD is_chirpy_red boolean NOT NULL default false;

-- +goose Down
ALTER TABLE users
DROP COLUMN is_chirpy_red;