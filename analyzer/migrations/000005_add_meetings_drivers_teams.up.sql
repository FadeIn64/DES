-- +goose Up
CREATE TABLE meetings (
                          meeting_key SERIAL PRIMARY KEY,
                          name TEXT NOT NULL,
                          description text NOT NULL,
                          circuit TEXT NOT NULL,
                          location TEXT NOT NULL,
                          start_date TIMESTAMPTZ NOT NULL,
                          end_date TIMESTAMPTZ NOT NULL,
                          year INTEGER NOT NULL,
                          dashboard_link TEXT
);

CREATE TABLE teams (
                       team_key SERIAL PRIMARY KEY,
                       name TEXT NOT NULL UNIQUE,
                       description text NOT NULL,
                       country TEXT,
                       color TEXT
);

CREATE TABLE drivers (
                         driver_number INTEGER PRIMARY KEY,
                         team_key INTEGER,
                         full_name TEXT NOT NULL,
                         abbreviation TEXT NOT NULL,
                         country TEXT,
                         date_of_birth DATE,
                         description text NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS drivers;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS meetings;