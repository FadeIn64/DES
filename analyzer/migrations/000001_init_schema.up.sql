-- +goose Up
CREATE TABLE laps (
                      meeting_key INTEGER NOT NULL,
                      session_key INTEGER NOT NULL,
                      driver_number INTEGER NOT NULL,
                      date_start TIMESTAMPTZ NOT NULL,
                      lap_duration DOUBLE PRECISION NOT NULL DEFAULT 0,
                      lap_number INTEGER NOT NULL,
                      sector_duration DOUBLE PRECISION[] NOT NULL,
                      info_time TIMESTAMPTZ NOT NULL,
                      is_pit_out_lap BOOLEAN NOT NULL DEFAULT FALSE,
                      updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                      PRIMARY KEY (driver_number, lap_number)
);

CREATE TABLE complete_laps (
                               meeting_key INTEGER NOT NULL,
                               session_key INTEGER NOT NULL,
                               driver_number INTEGER NOT NULL,
                               date_start TIMESTAMPTZ NOT NULL,
                               lap_duration DOUBLE PRECISION NOT NULL DEFAULT 0,
                               lap_number INTEGER NOT NULL,
                               sector_duration DOUBLE PRECISION[] NOT NULL,
                               info_time TIMESTAMPTZ NOT NULL,
                               is_pit_out_lap BOOLEAN NOT NULL DEFAULT FALSE,
                               updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                               PRIMARY KEY (driver_number, lap_number)
);

-- +goose Down
DROP TABLE IF EXISTS complete_laps;
DROP TABLE IF EXISTS laps;