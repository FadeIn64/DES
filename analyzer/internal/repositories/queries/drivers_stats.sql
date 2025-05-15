-- name: GetDriversStats :many
select * from drivers_stats_with_positions
where meeting_key = $1 and session_key = $2
order by position;

-- name: GetDriverStats :one
select * from drivers_stats_with_positions
where meeting_key = $1 and session_key = $2 and driver_number = $3;

-- name: GetDriverByPosition :one
select * from drivers_stats_with_positions
where meeting_key = $1 and session_key = $2 and position = $3;

-- name: UpsertDriversInterval :exec
insert into drivers_intervals (meeting_key, session_key, driver_number, interval, prediction_laps_to_overtake)
values ($1, $2, $3, $4, $5)
on conflict (meeting_key, session_key, driver_number) do update
set
    interval = excluded.interval,
    prediction_laps_to_overtake = excluded.prediction_laps_to_overtake;

