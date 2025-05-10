-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION get_driver_segment_pace(
    p_driver_number INTEGER,
    p_current_lap INTEGER
) RETURNS TABLE(average_pace FLOAT8, lap_count INTEGER) AS $$
BEGIN
RETURN QUERY
    WITH segment AS (
        SELECT lap_number, lap_duration
        FROM laps
        WHERE driver_number = p_driver_number
          AND lap_number <= p_current_lap
          AND is_pit_out_lap = false
          AND lap_duration > 0
        ORDER BY lap_number DESC
        LIMIT (
            SELECT COALESCE(
                (SELECT MIN(lap_number)
                 FROM laps
                 WHERE driver_number = p_driver_number
                   AND lap_number <= p_current_lap
                   AND is_pit_out_lap = true),
                p_current_lap
            )
        )
    )
    SELECT
        AVG(lap_duration)::FLOAT8,
            COUNT(*)::INTEGER
    FROM segment;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd


-- +goose Down
DROP FUNCTION IF EXISTS get_driver_segment_pace;