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
                               LIKE laps INCLUDING ALL,
                               PRIMARY KEY (driver_number, lap_number)
);

CREATE OR REPLACE FUNCTION move_complete_lap()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.lap_duration > 0 THEN
        INSERT INTO complete_laps VALUES (NEW.*)
        ON CONFLICT (driver_number, lap_number)
        DO UPDATE SET
    meeting_key = EXCLUDED.meeting_key,
                      session_key = EXCLUDED.session_key,
                      date_start = EXCLUDED.date_start,
                      lap_duration = EXCLUDED.lap_duration,
                      sector_duration = EXCLUDED.sector_duration,
                      info_time = EXCLUDED.info_time,
                      is_pit_out_lap = EXCLUDED.is_pit_out_lap,
                      updated_at = EXCLUDED.updated_at;
END IF;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_move_complete_lap
    AFTER INSERT OR UPDATE ON laps
                        FOR EACH ROW EXECUTE FUNCTION move_complete_lap();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_move_complete_lap ON laps;
DROP FUNCTION IF EXISTS move_complete_lap;
DROP TABLE IF EXISTS complete_laps;
DROP TABLE IF EXISTS laps;