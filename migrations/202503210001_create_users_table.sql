-- +goose Up
CREATE TABLE users
(
    id          VARCHAR(255) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_verified BOOLEAN      NOT NULL DEFAULT FALSE,
    verified_at TIMESTAMP NULL DEFAULT NULL
);


-- +goose Down
DROP TABLE IF EXISTS users;
