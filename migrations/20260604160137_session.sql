-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS session (
                                     refresh_token UUID PRIMARY KEY,
                                     user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
                                        is_active BOOLEAN NOT NULL DEFAULT true,
                                        created_at TIMESTAMP DEFAULT NOW(),
                                        expires_at TIMESTAMP NOT NULL
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE session;
-- +goose StatementEnd
