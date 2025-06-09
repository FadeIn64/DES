-- name: GetDriverByNumber :one
SELECT * FROM drivers WHERE driver_number = $1;

-- name: GetDriversByTeam :many
SELECT * FROM drivers WHERE team_key = $1 ORDER BY driver_number;

-- name: UpsertDrivers :exec
INSERT INTO drivers(driver_number, team_key, full_name, abbreviation, country, date_of_birth, description)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (driver_number)
DO UPDATE
SET
    team_key = excluded.team_key,
    full_name = excluded.full_name,
    abbreviation = excluded.abbreviation,
    country = excluded.country,
    date_of_birth = excluded.date_of_birth,
    description = excluded.description;