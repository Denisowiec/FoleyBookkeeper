-- name: CreateCalculation :one
INSERT INTO calculations (
    project_id,
    budget,
    currency,
    exchange_rate
) VALUES (
    $1,
    $2,
    $3,
    $4
) RETURNING *;

-- name: UpdateCalculation :one
UPDATE calculations SET
    project_id = $2,
    budget = $3,
    currency = $4,
    exchange_rate = $5,
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: AddEpisodeToCalculation :one
INSERT INTO episode_calc (
    episode_id,
    calc_id
) VALUES (
    $1,
    $2
) RETURNING *;

-- name: RemoveEpisodeFromCalculation :one
DELETE FROM episode_calc WHERE calc_id = $1 AND episode_id = $2 RETURNING *;

-- name: GetEpisodesForCalculation :many
SELECT episode_id FROM episode_calc WHERE calc_id = $1;

-- name: GetCalculation :one
SELECT * FROM calculations WHERE id = $1;

-- name: GetAllCalculationsForProject :many
SELECT * FROM calculations WHERE project_id = $1;

-- name: GetMinutesForCalculation :one
SELECT SUM(sessions.duration) AS minutes FROM sessions
JOIN episode_calc ON episode_calc.episode_id = sessions.episode_id
WHERE episode_calc.calc_id = $1;
