-- name: GetDriversStats :many
select * from drivers_stats_with_positions
where meeting_key = $1 and session_key = $2
order by position;

-- name: GetDriverStats :one
select * from drivers_stats_with_positions
where meeting_key = $1 and session_key = $2 and driver_number = $3;

-- name: UpsertDriverStats :exec
insert into drivers_stats (meeting_key, session_key, driver_number, date_start, date_end, lap_duration, lap_number)
values ($1, $2, $3, $4, $5, $6, $7)
on conflict (meeting_key, session_key, driver_number)
do update set
              driver_number = excluded.driver_number,
              date_start = excluded.date_start,
              date_end = excluded.date_end,
              lap_duration = excluded.lap_duration,
              lap_number = excluded.lap_number;
