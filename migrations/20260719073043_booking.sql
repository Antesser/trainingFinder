-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS training_booking
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    training_id TEXT      NOT NULL,
    booked_by   UUID      NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    status      text,
    book_from   TIMESTAMP NOT NULL,
    book_to     TIMESTAMP NOT NULL,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE training_booking;
-- +goose StatementEnd
