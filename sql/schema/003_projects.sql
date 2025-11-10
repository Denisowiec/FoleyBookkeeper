-- +goose Up
CREATE TABLE projects (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL,
    title TEXT UNIQUE NOT NULL,
    client TEXT NOT NULL
);

-- +goose Down
DROP TABLE projects;