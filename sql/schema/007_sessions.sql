-- +goose Up
CREATE TYPE part AS ENUM ('props', 'footsteps', 'movements', 'dialogue', 'adr', 'music', 'background', 'other');
CREATE TYPE activity AS ENUM ('record', 'edit', 'service', 'spotting', 'other');

CREATE TABLE sessions (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    session_date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    episode_id UUID NOT NULL REFERENCES episodes ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES projects ON DELETE CASCADE,
    duration INTEGER NOT NULL,
    part_worked_on PART NOT NULL,
    activity_done ACTIVITY NOT NULL
);

-- +goose Down
DROP TABLE sessions;
DROP TYPE part;
DROP TYPE activity;