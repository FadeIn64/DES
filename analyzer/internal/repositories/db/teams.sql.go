// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: teams.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getAllTeams = `-- name: GetAllTeams :many
SELECT team_key, name, description, country, color FROM teams ORDER BY name
`

func (q *Queries) GetAllTeams(ctx context.Context) ([]Team, error) {
	rows, err := q.db.Query(ctx, getAllTeams)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Team
	for rows.Next() {
		var i Team
		if err := rows.Scan(
			&i.TeamKey,
			&i.Name,
			&i.Description,
			&i.Country,
			&i.Color,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTeamByID = `-- name: GetTeamByID :one
SELECT team_key, name, description, country, color FROM teams WHERE team_key = $1
`

func (q *Queries) GetTeamByID(ctx context.Context, teamKey int64) (Team, error) {
	row := q.db.QueryRow(ctx, getTeamByID, teamKey)
	var i Team
	err := row.Scan(
		&i.TeamKey,
		&i.Name,
		&i.Description,
		&i.Country,
		&i.Color,
	)
	return i, err
}

const upsertTeam = `-- name: UpsertTeam :exec
INSERT INTO teams(team_key, name, description, country, color)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (team_key)
DO UPDATE
SET
    name = excluded.name,
    description = excluded.description,
    country = excluded.country,
    color = excluded.color
`

type UpsertTeamParams struct {
	TeamKey     int64
	Name        string
	Description string
	Country     pgtype.Text
	Color       pgtype.Text
}

func (q *Queries) UpsertTeam(ctx context.Context, arg UpsertTeamParams) error {
	_, err := q.db.Exec(ctx, upsertTeam,
		arg.TeamKey,
		arg.Name,
		arg.Description,
		arg.Country,
		arg.Color,
	)
	return err
}
