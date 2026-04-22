-- +goose Up
CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    user_name  VARCHAR   NOT NULL,
    email      VARCHAR   NOT NULL,
    password   VARCHAR   NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE users;
