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

func (r *LapRepository) ProcessLap(ctx context.Context, lap models.Lap) (*models.LapAnalysis, error) {
	var analysis models.LapAnalysis

	err := r.manager.Do(ctx, func(ctx context.Context) error {
		q := db.New(trmpgx.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db))

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

		// 2. Если это пит-стоп, завершаем обработку
		if lap.IsPitOutLap {
			return nil
		}

		analysis.DriverNumber = lap.DriverNumber
		analysis.CurrentLapTime = lap.LapDuration

		// 3. Рассчитываем среднее время круга (исключая пит-стопы)
		analysis.AverageLapTime, err = q.GetAverageLapTime(ctx, db.GetAverageLapTimeParams{
			DriverNumber: lap.DriverNumber,
			IsPitOutLap:  false,
		})
		if err != nil {
			return fmt.Errorf("get average lap time: %w", err)
		}

		// 4. Рассчитываем средний темп на текущем сегменте
		segment, err := q.GetCurrentSegmentPace(ctx, db.GetCurrentSegmentPaceParams{
			DriverNumber: lap.DriverNumber,
			LapNumber:    lap.LapNumber,
		})
		if err != nil {
			return fmt.Errorf("get segment pace: %w", err)
		}

		analysis.AverageSegmentPace = segment.AveragePace
		analysis.LapsInSegment = int(segment.LapCount)

		// 5. Сравнение текущего круга со средним
		if analysis.AverageLapTime > 0 {
			analysis.ComparisonWithAvg = (analysis.CurrentLapTime - analysis.AverageLapTime) / analysis.AverageLapTime * 100
		}

		// 6. Определяем тренд
		analysis.PositionTrend = r.calculateTrend(
			analysis.CurrentLapTime,
			analysis.AverageSegmentPace,
		)

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &analysis, nil
}

func (r *LapRepository) calculateTrend(current, average float64) string {
	if average == 0 {
		return "stable"
	}

	ratio := (current - average) / average
	switch {
	case ratio < -0.03: // >3% лучше
		return "improving"
	case ratio > 0.03: // >3% хуже
		return "declining"
	default:
		return "stable"
	}
}
