-- name: GetCurrentMeeting :one
SELECT * FROM meetings
WHERE start_date <= NOW() AND end_date >= NOW()
LIMIT 1;

-- name: GetMeetingByKey :one
SELECT * from meetings
    WHERE meeting_key = $1;

-- name: GetMeetings :many
SELECT * FROM meetings
    ORDER BY start_date;

-- name: UpsertMeeting :exec
INSERT INTO meetings (meeting_key, name, description, circuit, location, start_date, end_date, year, dashboard_link) VALUES (
             $1, $2, $3, $4, $5, $6, $7, $8, $9
         )
ON CONFLICT (meeting_key)
do update
set
    name = excluded.name,
    description = excluded.description,
    circuit = excluded.circuit,
    location = excluded.location,
    start_date = excluded.start_date,
    end_date = excluded.end_date,
    year = excluded.year,
    dashboard_link = excluded.dashboard_link;
