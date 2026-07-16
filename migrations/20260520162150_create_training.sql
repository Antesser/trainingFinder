-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS training (
                                        id TEXT,
                                        trainer_id TEXT NOT NULL,
                                        user_id TEXT NOT NULL,
                                        started_at TIMESTAMP  NOT NULL,
                                        ended_at TIMESTAMP NOT NULL,
                                        booked_by UUID,
                                        additional_info TEXT NOT NULL,
                                        PRIMARY KEY(id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE training;
-- +goose StatementEnd
