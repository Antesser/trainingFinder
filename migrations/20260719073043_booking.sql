-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS booking (
                                        id TEXT,
                                        training_id TEXT NOT NULL,
                                        booked_by UUID NOT NULL,
                                        time_booked TIMESTAMP  NOT NULL,
                                        PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE booking;
-- +goose StatementEnd
