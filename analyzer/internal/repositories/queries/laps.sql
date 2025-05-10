-- name: GetLap :one
SELECT * FROM laps
WHERE driver_number = $1 AND lap_number = $2;

-- name: GetAverageLapTime :one
SELECT AVG(lap_duration)::float8
FROM laps
WHERE driver_number = $1
  AND is_pit_out_lap = $2
  AND lap_duration > 0;

-- name: GetCurrentSegmentPace :one
WITH segment AS (
    SELECT lap_number, lap_duration
    FROM laps l
    WHERE l.driver_number = $1
      AND l.lap_number <= $2
      AND l.is_pit_out_lap = false
      AND l.lap_duration > 0
    ORDER BY l.lap_number DESC
    LIMIT (
        SELECT COALESCE(
                       (SELECT MIN(l2.lap_number)
                        FROM laps l2
                        WHERE l2.driver_number = $1
                          AND l2.lap_number <= $2
                          AND l2.is_pit_out_lap = true),
                       $2
               )
        )
    )
    SELECT
        AVG(lap_duration)::float8 as average_pace,
        COUNT(*) as lap_count
    FROM segment;

-- name: UpsertLap :exec
INSERT INTO laps (
    meeting_key, session_key, driver_number,
    date_start, lap_duration, lap_number,
    sector_duration, info_time, is_pit_out_lap
) VALUES (
             $1, $2, $3, $4, $5, $6, $7, $8, $9
         )
ON CONFLICT (driver_number, lap_number)
    DO UPDATE SET
                  meeting_key = EXCLUDED.meeting_key,
                  session_key = EXCLUDED.session_key,
                  date_start = EXCLUDED.date_start,
                  lap_duration = COALESCE(NULLIF(EXCLUDED.lap_duration, 0), laps.lap_duration),
                  sector_duration = EXCLUDED.sector_duration,
                  info_time = GREATEST(EXCLUDED.info_time, laps.info_time),
                  is_pit_out_lap = EXCLUDED.is_pit_out_lap,
                  updated_at = NOW();

-- name: MoveCompleteLap :exec
INSERT INTO complete_laps
SELECT * FROM laps l
WHERE l.driver_number = $1 AND l.lap_number = $2 AND l.lap_duration > 0
ON CONFLICT (driver_number, lap_number) DO NOTHING;