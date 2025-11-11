-- +goose Up
CREATE TABLE projects (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    title TEXT UNIQUE NOT NULL,
    client_id UUID NOT NULL REFERENCES clients ON DELETE CASCADE
);

-- +goose Down
DROP TABLE projects;