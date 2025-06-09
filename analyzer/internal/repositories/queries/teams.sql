-- name: GetTeamByID :one
SELECT * FROM teams WHERE team_key = $1;

-- name: GetAllTeams :many
SELECT * FROM teams ORDER BY name;

-- name: UpsertTeam :exec
INSERT INTO teams(team_key, name, description, country, color)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (team_key)
DO UPDATE
SET
    name = excluded.name,
    description = excluded.description,
    country = excluded.country,
    color = excluded.color;