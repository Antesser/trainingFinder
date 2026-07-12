-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS roles (
                                     id BIGSERIAL PRIMARY KEY,
                                     role text NOT NULL,
                                     PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
Drop table roles;
-- +goose StatementEnd
