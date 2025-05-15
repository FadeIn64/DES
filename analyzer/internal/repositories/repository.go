package repositories

import (
	"DAS/internal/repositories/db"
	"DAS/models"
	"context"
	"fmt"
	trmpgx "github.com/avito-tech/go-transaction-manager/pgxv5"
	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"math"
	"time"
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

const MaxLapsToOvertake = 11

func (r *LapRepository) ProcessLap(ctx context.Context, lap models.Lap) (*models.LapAnalysis, error) {
	var analysis models.LapAnalysis

	err := r.manager.Do(ctx, func(ctx context.Context) error {
		var err error
		q := db.New(trmpgx.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db))

		// Считаем дополнительные параметры
		lapDuration := float64(0)
		completedSectors := 0
		if lap.LapDuration != 0 {
			lapDuration = lap.LapDuration
			completedSectors = len(lap.SectorDuration)
		} else {
			for _, sector := range lap.SectorDuration {
				lapDuration += sector
				if sector > 0 {
					completedSectors++
				}
			}
		}
		timeEnd := int64(lapDuration * 1000)
		dateEnd := lap.DateStart.Add(time.Millisecond * time.Duration(timeEnd))

		//. Вставляем/обновляем запись
		if err := q.UpsertLap(ctx, db.UpsertLapParams{
			MeetingKey:       lap.MeetingKey,
			SessionKey:       lap.SessionKey,
			DriverNumber:     lap.DriverNumber,
			CompletedSectors: int32(completedSectors),
			DateStart:        pgtype.Timestamptz{Time: lap.DateStart, Valid: true},
			DateEnd:          pgtype.Timestamptz{Time: dateEnd, Valid: true},
			LapDuration:      lap.LapDuration,
			LapNumber:        lap.LapNumber,
			SectorDuration:   lap.SectorDuration,
			InfoTime:         pgtype.Timestamptz{Time: lap.InfoTime, Valid: true},
			IsPitOutLap:      lap.IsPitOutLap,
		}); err != nil {
			return fmt.Errorf("upsert lap: %w", err)
		}

		// Если круг завершен, перемещаем его
		if lap.LapDuration > 0 {
			if err := q.MoveCompleteLap(ctx, db.MoveCompleteLapParams{
				DriverNumber:     lap.DriverNumber,
				LapNumber:        lap.LapNumber,
				MeetingKey:       lap.MeetingKey,
				SessionKey:       lap.SessionKey,
				CompletedSectors: int32(len(lap.SectorDuration)),
			}); err != nil {
				return fmt.Errorf("move complete lap: %w", err)
			}
		}

		//. Если это пит-стоп, завершаем обработку
		if lap.IsPitOutLap {
			return nil
		}

		analysis.DriverNumber = lap.DriverNumber
		analysis.CurrentLapTime = lap.LapDuration
		analysis.LapNumber = lap.LapNumber
		analysis.MeetingKey = lap.MeetingKey
		analysis.SessionKey = lap.SessionKey

		// Рассчитываем среднее время круга (исключая пит-стопы)
		analysis.AverageLapTime, err = q.GetAverageLapTime(ctx, db.GetAverageLapTimeParams{
			DriverNumber: lap.DriverNumber,
			IsPitOutLap:  false,
			MeetingKey:   lap.MeetingKey,
			SessionKey:   lap.SessionKey,
		})
		if err != nil {
			return fmt.Errorf("get average lap time: %w", err)
		}

		// Рассчитываем средний темп на текущем сегменте
		segment, err := q.GetCurrentSegmentPace(ctx, db.GetCurrentSegmentPaceParams{
			DriverNumber: lap.DriverNumber,
			LapNumber:    lap.LapNumber,
			MeetingKey:   lap.MeetingKey,
			SessionKey:   lap.SessionKey,
		})
		if err != nil {
			return fmt.Errorf("get segment pace: %w", err)
		}

		analysis.AverageSegmentPace = segment.AveragePace
		analysis.LapsInSegment = int(segment.LapCount)

		// Сравнение текущего круга со средним
		if analysis.AverageLapTime > 0 {
			analysis.ComparisonWithAvg = (analysis.CurrentLapTime - analysis.AverageLapTime) / analysis.AverageLapTime * 100
		}

		// Определяем тренд
		analysis.PositionTrend = r.calculateTrend(
			analysis.CurrentLapTime,
			analysis.AverageSegmentPace,
		)

		// Определяем интервалы от следующего гонщика
		curDriver, err := q.GetDriverStats(ctx, db.GetDriverStatsParams{
			DriverNumber: lap.DriverNumber,
			MeetingKey:   lap.MeetingKey,
			SessionKey:   lap.SessionKey,
		})
		if err != nil {
			return fmt.Errorf("get cur driver stats: %w", err)
		}

		// Если гонщик едет первым, то считать отставние не откого
		if curDriver.Position == 1 {
			err = q.UpsertDriversInterval(ctx, db.UpsertDriversIntervalParams{
				MeetingKey:               lap.MeetingKey,
				SessionKey:               lap.SessionKey,
				DriverNumber:             lap.DriverNumber,
				Interval:                 0,
				PredictionLapsToOvertake: pgtype.Int4{Int32: MaxLapsToOvertake, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("upsert driver interval: %w", err)
			}
			return nil
		}

		nextDriver, err := q.GetDriverByPosition(ctx, db.GetDriverByPositionParams{
			MeetingKey: lap.MeetingKey,
			SessionKey: lap.SessionKey,
			Position:   curDriver.Position - 1,
		})

		if err != nil {
			return fmt.Errorf("get next driver: %w", err)
		}

		nextDriverLap, err := q.GetLap(ctx, db.GetLapParams{
			DriverNumber:     nextDriver.DriverNumber,
			LapNumber:        lap.LapNumber,
			MeetingKey:       lap.MeetingKey,
			SessionKey:       lap.SessionKey,
			CompletedSectors: int32(completedSectors),
		})
		if err != nil {
			return fmt.Errorf("get next driver lap: %w", err)
		}

		timeInterval := nextDriverLap.DateEnd.Time.Sub(dateEnd)
		interval := float64(timeInterval) / float64(time.Second)

		predictOvertake := MaxLapsToOvertake
		nextDriverSegmentPace, err := q.GetCurrentSegmentPace(ctx, db.GetCurrentSegmentPaceParams{
			DriverNumber: nextDriver.DriverNumber,
			MeetingKey:   lap.MeetingKey,
			SessionKey:   lap.SessionKey,
			LapNumber:    lap.LapNumber,
		})
		if err != nil {
			return fmt.Errorf("get next driver segment pace: %w", err)
		}

		if analysis.AverageSegmentPace != 0 && nextDriverSegmentPace.AveragePace != 0 {
			paceDiff := analysis.AverageSegmentPace - nextDriverSegmentPace.AveragePace
			if paceDiff > 0 {
				predictOvertakeFloat := -1 * math.Ceil(interval/paceDiff)
				if predictOvertakeFloat < 11 {
					predictOvertake = int(predictOvertakeFloat)
				}
			}
		}

		err = q.UpsertDriversInterval(ctx, db.UpsertDriversIntervalParams{
			MeetingKey:               lap.MeetingKey,
			SessionKey:               lap.SessionKey,
			DriverNumber:             lap.DriverNumber,
			Interval:                 interval,
			PredictionLapsToOvertake: pgtype.Int4{Int32: int32(predictOvertake), Valid: true},
		})
		if err != nil {
			return fmt.Errorf("upsert driver interval: %w", err)
		}

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
