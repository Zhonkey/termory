-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID NOT NULL PRIMARY KEY,
    role VARCHAR(50) NOT NULL ,
    email VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
