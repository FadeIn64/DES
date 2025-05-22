-- name: GetDriverByNumber :one
SELECT * FROM drivers WHERE driver_number = $1;

-- name: GetDriversByTeam :many
SELECT * FROM drivers WHERE team_key = $1 ORDER BY driver_number;