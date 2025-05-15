-- +goose Up
create table drivers_intervals(
    meeting_key INTEGER NOT NULL,
    session_key INTEGER NOT NULL,
    driver_number INTEGER NOT NULL,
    interval DOUBLE PRECISION NOT NULL DEFAULT 0,
    prediction_laps_to_overtake INTEGER,

    PRIMARY KEY (meeting_key, session_key, driver_number)
);


-- +goose Down
DROP TABLE drivers_intervals