-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS roles (
                                     id text PRIMARY KEY,
                                     role text NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
Drop table roles;
-- +goose StatementEnd
