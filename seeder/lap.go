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

	laps, err := readers.ReadLapsData()
	if err != nil {
		return err
	}

	q := db.New(pool)

	for _, lap := range laps {
		args := convertLapToInsert(&lap)
		err := q.InsertLap(ctx, *args)
		if err != nil {
			return err
		}
	}

	return nil
}

func convertLapToInsert(lap *models.F1Lap) *db.InsertLapParams {
	return &db.InsertLapParams{
		MeetingKey:      lap.MeetingKey,
		SessionKey:      lap.SessionKey,
		DriverNumber:    pgtype.Int4{Int32: lap.DriverNumber, Valid: true},
		I1Speed:         pgtype.Int4{Int32: lap.I1Speed, Valid: true},
		I2Speed:         pgtype.Int4{Int32: lap.I2Speed, Valid: true},
		StSpeed:         pgtype.Int4{Int32: lap.StSpeed, Valid: true},
		DateStart:       pgtype.Timestamp{Time: lap.DateStart, Valid: true},
		LapDuration:     pgtype.Float8{Float64: lap.LapDuration, Valid: true},
		IsPitOutLap:     pgtype.Bool{Bool: lap.IsPitOutLap, Valid: true},
		DurationSector1: pgtype.Float8{Float64: lap.DurationSector1, Valid: true},
		DurationSector2: pgtype.Float8{Float64: lap.DurationSector2, Valid: true},
		DurationSector3: pgtype.Float8{Float64: lap.DurationSector3, Valid: true},
		SegmentsSector1: lap.SegmentsSector1,
		SegmentsSector2: lap.SegmentsSector2,
		SegmentsSector3: lap.SegmentsSector3,
		LapNumber:       lap.LapNumber,
	}
}
