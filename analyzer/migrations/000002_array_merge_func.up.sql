-- +goose Up
CREATE OR REPLACE FUNCTION array_merge_non_null(
    arr1 DOUBLE PRECISION[],
    arr2 DOUBLE PRECISION[]
) RETURNS DOUBLE PRECISION[] AS $$
DECLARE
result DOUBLE PRECISION[];
    i INTEGER;
BEGIN
    result = arr1;
FOR i IN 1..LEAST(array_length(arr1, 1), array_length(arr2, 1)) LOOP
        IF arr2[i] IS NOT NULL THEN
            result[i] = arr2[i];
END IF;
END LOOP;
RETURN result;
END;
$$ LANGUAGE plpgsql;

-- +goose Down
DROP FUNCTION IF EXISTS array_merge_non_null;