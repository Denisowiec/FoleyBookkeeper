-- +goose Up
CREATE TABLE user_session (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    session_id UUID NOT NULL REFERENCES sessions ON DELETE CASCADE
);

-- +goose Down
DROP TABLE user_session;