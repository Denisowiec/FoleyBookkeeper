-- name: CreateSession :one
INSERT INTO sessions (
    duration,
    episode_id,
    project_id,
    part_worked_on,
    activity_done
) VALUES (
    $1,
    $2,
    (SELECT project_id FROM episodes WHERE id = $2),
    $3,
    $4
) RETURNING * ;

-- name: UpdateSession :one
UPDATE sessions SET
    duration = $2,
    episode_id =$3,
    project_id = (SELECT project_id FROM episodes WHERE id = $3),
    part_worked_on = $4,
    activity_done = $5
WHERE sessions.id = $1 RETURNING *;

-- name: GetSession :one
SELECT 
    sessions.id,
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


-- name: GetAllSessionsForProject :many
SELECT * FROM sessions WHERE project_id = $1;

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
