package main

import (
	"DES/db"
	"DES/models"
	"DES/readers"
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func seedLaps(ctx context.Context, pool *pgxpool.Pool) error {

	f1laps, err := readers.ReadLapsData()
	if err != nil {
		return err
	}

	q := db.New(pool)

	laps := make([]models.Lap, 0, len(f1laps)*3)

	for _, f1lap := range f1laps {
		laps = append(laps, models.F1LapToSectorLapsWithTiming(f1lap)...)
	}

	for _, lap := range laps {
		args := convertLapToInsert(&lap)
		err := q.InsertLap(ctx, *args)
		if err != nil {
			return err
		}
	}

	return nil
}

func convertLapToInsert(lap *models.Lap) *db.InsertLapParams {
	return &db.InsertLapParams{
		MeetingKey:     lap.MeetingKey,
		SessionKey:     lap.SessionKey,
		DriverNumber:   lap.DriverNumber,
		DateStart:      pgtype.Timestamptz{Time: lap.DateStart, Valid: true},
		LapDuration:    lap.LapDuration,
		LapNumber:      lap.LapNumber,
		SectorDuration: lap.SectorDuration,
		InfoTime:       pgtype.Timestamptz{Time: lap.InfoTime, Valid: true},
	}
}
