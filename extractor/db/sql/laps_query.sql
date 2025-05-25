-- name: InsertLap :exec
insert into laps(meeting_key, session_key, driver_number, date_start, lap_duration, lap_number, sector_duration, info_time, is_pit_out_lap)
values ($1, $2, $3,$4, $5, $6,$7, $8, $9);
-- name: GetLapsStartDateBetween :many
select *
from laps
where info_time >= $1
  and info_time <= $2
order by info_time;