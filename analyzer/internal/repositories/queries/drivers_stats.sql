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

