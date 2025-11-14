-- name: CreateEpisode :one
INSERT INTO episodes (
    title,
    episode_number,
    project_id
) VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: UpdateEpisode :one
UPDATE episodes SET
    title = $2,
    episode_number = $3,
    project_id = $4,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: GetEpisodeByID :one
SELECT * FROM episodes WHERE id = $1;

-- name: GetEpisodeByNumber :one
SELECT * FROM episodes WHERE project_id = $1 AND episode_number = $2;

-- name: GetAllEpisodes :many
SELECT * FROM episodes;

-- name: GetAllEpisodesForProject :many
SELECT * FROM episodes WHERE project_id = $1 ORDER BY episode_number ASC;

