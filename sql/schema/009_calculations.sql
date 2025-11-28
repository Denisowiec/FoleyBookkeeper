-- +goose Up
CREATE TABLE calculations (
    id UUID DEFAULT gen_random_uuid() PRiMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    project_id NOT NULL REFERENCES projects ON DELETE CASCADE,
    budget NUMERIC NOT NULL DEFAULT 0,
    currency TEXT NOT NULL DEFAULT "PLN",
    exchange_rate NUMERIC NOT NULL DEFAULT 1

);

-- +goose Down
DROP TABLE calculations;