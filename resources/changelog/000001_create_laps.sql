-- +goose Up
-- +goose StatementBegin
create sequence if not exists laps_seq
    increment by 1;

create table if not exists laps(
    id  bigint default  nextval('laps_seq'::regclass) not null primary key,
    meeting_key INT NOT NULL,
    session_key INT NOT NULL,
    driver_number INT,
    i1_speed INT,
    i2_speed INT,
    st_speed INT,
    date_start timestamp,
    lap_duration FLOAT,
    is_pit_out_lap BOOLEAN,
    duration_sector_1 FLOAT,
    duration_sector_2 FLOAT,
    duration_sector_3 FLOAT,
    segments_sector_1 INT[],
    segments_sector_2 INT[],
    segments_sector_3 INT[],
    lap_number INT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table laps;
drop sequence laps_seq;

-- +goose StatementEnd
