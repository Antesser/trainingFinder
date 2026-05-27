-- +goose Up
-- +goose StatementBegin
CREATE TABLE outbox
(
    id BIGSERIAL PRIMARY KEY,
    channel  TEXT NOT NULL,
    value TEXT NOT NULL,
    message_value  TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE outbox
-- +goose StatementEnd
