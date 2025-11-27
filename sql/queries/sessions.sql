-- name: CreateSession :one
INSERT INTO sessions (
    duration,
    session_date,
    episode_id,
    project_id,
    part_worked_on,
    activity_done
) VALUES (
    $1,
    $2,
    $3,
    (SELECT project_id FROM episodes WHERE id = $3),
    $4,
    $5
) RETURNING * ;

-- name: UpdateSession :one
UPDATE sessions SET
    duration = $2,
    session_date = $3,
    episode_id = $4,
    project_id = (SELECT project_id FROM episodes WHERE id = $4),
    part_worked_on = $5,
    activity_done = $6
WHERE sessions.id = $1 RETURNING *;

-- name: GetSession :one
SELECT 
    sessions.id,
    sessions.session_date,
    sessions.created_at,
    sessions.updated_at,
    sessions.duration,
    sessions.part_worked_on,
    sessions.activity_done,
    episodes.id AS episode_id,
    episodes.title AS episode_title,
    episodes.episode_number AS episode_number,
    projects.id AS project_id,
    projects.title AS project_title
FROM sessions
JOIN episodes ON episodes.id = sessions.episode_id
JOIN projects ON projects.id = sessions.project_id WHERE sessions.id = $1;

-- name: GetSessionsForProject :many
SELECT * FROM sessions WHERE project_id = $1 ORDER BY session_date DESC LIMIT $2;

-- name: GetSessionsForEpisode :many
SELECT * FROM sessions WHERE episode_id = $1 ORDER BY session_date DESC LIMIT $2;

-- name: GetSessions :many
SELECT * FROM sessions ORDER BY session_date DESC LIMIT $1;

-- name: DeleteSession :one
DELETE FROM sessions WHERE id = $1 RETURNING *;

-- name: AddUserToSession :one
INSERT INTO user_session (
    user_id,
    session_id
) VALUES (
    $1,
    $2
) RETURNING *;

-- name: GetUsersForSession :many
SELECT user_id FROM user_session WHERE session_id = $1;
