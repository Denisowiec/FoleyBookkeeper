-- +goose Up
CREATE TABLE episodes (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL,
    title TEXT,
    project_id UUID NOT NULL REFERENCES projects ON DELETE CASCADE
);

-- +goose Down
DROP TABLE episodes;