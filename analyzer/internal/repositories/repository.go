package repositories

import (
	"DAS/internal/repositories/db"
	"DAS/models"
	"context"
	"errors"
	"fmt"
	trmpgx "github.com/avito-tech/go-transaction-manager/pgxv5"
	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type LapRepository struct {
	db      *pgxpool.Pool
	manager trm.Manager
}

func NewLapRepository(db *pgxpool.Pool, manager trm.Manager) *LapRepository {
	return &LapRepository{
		db:      db,
		manager: manager,
	}
}

func (r *LapRepository) ProcessLap(ctx context.Context, lap models.Lap) error {
	return r.manager.Do(ctx, func(ctx context.Context) error {
		q := db.New(trmpgx.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db))

		log.Printf("processing lap %+v", lap)

		// 1. Проверяем существование записи
		_, err := q.GetLap(ctx, db.GetLapParams{
			DriverNumber: lap.DriverNumber,
			LapNumber:    lap.LapNumber,
		})
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("get lap: %w", err)
		}

		// 2. Вставляем/обновляем запись
		if err := q.UpsertLap(ctx, db.UpsertLapParams{
			MeetingKey:     lap.MeetingKey,
			SessionKey:     lap.SessionKey,
			DriverNumber:   lap.DriverNumber,
			DateStart:      pgtype.Timestamptz{Time: lap.DateStart, Valid: true},
			LapDuration:    lap.LapDuration,
			LapNumber:      lap.LapNumber,
			SectorDuration: lap.SectorDuration,
			InfoTime:       pgtype.Timestamptz{Time: lap.InfoTime, Valid: true},
			IsPitOutLap:    lap.IsPitOutLap,
		}); err != nil {
			return fmt.Errorf("upsert lap: %w", err)
		}

		// 3. Если круг завершен, перемещаем его
		if lap.LapDuration > 0 {
			if err := q.MoveCompleteLap(ctx, db.MoveCompleteLapParams{
				DriverNumber: lap.DriverNumber,
				LapNumber:    lap.LapNumber,
			}); err != nil {
				return fmt.Errorf("move complete lap: %w", err)
			}
		}

		return nil
	})
}
