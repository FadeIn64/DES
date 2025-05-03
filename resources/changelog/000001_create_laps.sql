-- +goose Up
-- +goose StatementBegin
create sequence if not exists laps_seq
    increment by 1;

create table if not exists laps(
    id  bigint default  nextval('laps_seq'::regclass) not null primary key,
    meeting_key integer not null,
    session_key integer not null,
    driver_number integer not null,
    date_start timestamptz not null,
    lap_duration double precision not null,
    lap_number integer not null,
    sector_duration double precision[3] not null,
    info_time timestamptz not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table laps;
drop sequence laps_seq;

-- +goose StatementEnd
