-- +goose Up
create table drivers_stats(
    meeting_key INTEGER NOT NULL,
    session_key INTEGER NOT NULL,
    driver_number INTEGER NOT NULL,
    date_start TIMESTAMPTZ NOT NULL,
    date_end TIMESTAMPTZ NOT NULL,
    lap_duration DOUBLE PRECISION NOT NULL DEFAULT 0,
    sectors INT NOT NULL,
    lap_number INTEGER NOT NULL,
    PRIMARY KEY (meeting_key, session_key, driver_number)
);

create view drivers_stats_with_positions as
    select row_number() over (order by lap_number desc, sectors desc, date_end) as position, * from drivers_stats;

-- +goose Down
drop view drivers_stats_with_positions;
drop table drivers_stats;