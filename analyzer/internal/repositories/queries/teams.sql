-- name: GetTeamByID :one
SELECT * FROM teams WHERE team_key = $1;

-- name: GetAllTeams :many
SELECT * FROM teams ORDER BY name;