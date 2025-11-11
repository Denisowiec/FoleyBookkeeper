-- +goose Up
CREATE TABLE clients (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    client_name TEXT UNIQUE NOT NULL,
    email TEXT,
    notes TEXT
);

-- +goose Down
DROP TABLE clients;
