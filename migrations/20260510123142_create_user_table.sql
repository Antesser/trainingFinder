-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
                                     id      TEXT PRIMARY KEY,
                                     login    VARCHAR(20) NOT NULL,
                                     password VARCHAR(50) NOT NULL,
                                     role_id  text        NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd