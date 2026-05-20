-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id TEXT NOT NULL,
    username VARCHAR(20) NOT NULL,
    password VARCHAR(15) NOT NULL,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE users;
