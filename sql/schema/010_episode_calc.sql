-- +goose Up
CREATE TABLE episode_calc (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    episode_id UUID NOT NULL REFERENCES episodes ON DELETE CASCADE,
    calc_id UUID NOT NULL REFERENCES calculations ON DELETE CASCADE,
    UNIQUE (episode_id, calc_id)
);

-- +goose Down
DROP TABLE episode_calc;