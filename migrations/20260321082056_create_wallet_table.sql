-- +goose Up
CREATE TABLE wallets
(
    id          SERIAL PRIMARY KEY,
    uuid        BIGINT    NOT NULL,
    balance     BIGINT    NOT NULL DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

create index idx_wallet_id on wallets (uuid);

-- +goose Down
DROP TABLE wallets;