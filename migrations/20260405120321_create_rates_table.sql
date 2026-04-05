-- +goose Up
CREATE TABLE rates (
    id         BIGSERIAL PRIMARY KEY,
    ask        BIGINT      NOT NULL,
    bid        BIGINT      NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE rates;
