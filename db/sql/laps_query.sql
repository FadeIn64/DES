-- name: InsertLap :exec
insert into laps (meeting_key, session_key, driver_number, i1_speed, i2_speed, st_speed, date_start, lap_duration, is_pit_out_lap, duration_sector_1, duration_sector_2, duration_sector_3, segments_sector_1, segments_sector_2, segments_sector_3, lap_number)
values
($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);

-- name: GetLapsStartDateBetween :many
select * from laps
where date_start >= $1 and date_start <= $2;