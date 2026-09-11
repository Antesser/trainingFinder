-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS training
(
    id              TEXT,
    trainer_id      TEXT      NOT NULL,
    available_from  TIMESTAMP NOT NULL,
    available_to    TIMESTAMP NOT NULL,
    duration        INTERVAL  NOT NULL,
    additional_info TEXT      NOT NULL,
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE training;
-- +goose StatementEnd
