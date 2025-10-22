-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_tokens (
    id UUID NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    expires_at TIMESTAMP,
    created_at TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE refresh_tokens;
-- +goose StatementEnd
