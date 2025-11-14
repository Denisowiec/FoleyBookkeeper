-- +goose Up
CREATE TABLE episodes (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    title TEXT,
    episode_number INTEGER NOT NULL,
    project_id UUID NOT NULL REFERENCES projects ON DELETE CASCADE,
    UNIQUE (episode_number, project_id)
);

-- +goose Down
DROP TABLE episodes;