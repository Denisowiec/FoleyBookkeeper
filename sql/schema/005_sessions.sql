-- +goose Up
CREATE TYPE part AS ENUM ('props', 'footsteps', 'movements', 'dialogue', 'adr', 'music', 'background');
CREATE TYPE activity AS ENUM ('record', 'edit', 'service', 'spotting');

CREATE TABLE sessions (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL,
    episode_id UUID NOT NULL REFERENCES episodes ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES projects ON DELETE CASCADE,
    duration INTERVAL NOT NULL,
    part_worked_on PART NOT NULL,
    activity_done ACTIVITY NOT NULL
);

-- +goose Down
DROP TABLE sessions;
DROP TYPE part;
DROP TYPE ativity;