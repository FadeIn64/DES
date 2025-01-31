package extractors

import (
	"DES/db"
	"DES/models"
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type LapExtractor interface {
	ExtractLap(ctx context.Context, from time.Time, to time.Time) ([]models.Lap, error)
}

func NewLapExtractor(pool *pgxpool.Pool) LapExtractor {
	return &lapExtractor{pool}
}

type lapExtractor struct {
	pool *pgxpool.Pool
}

func (l *lapExtractor) ExtractLap(ctx context.Context, from time.Time, to time.Time) ([]models.Lap, error) {

	q := db.New(l.pool)

	args := db.GetLapsStartDateBetweenParams{
		DateStart:   pgtype.Timestamp{Time: from, Valid: true},
		DateStart_2: pgtype.Timestamp{Time: to, Valid: true},
	}

	laps, err := q.GetLapsStartDateBetween(ctx, args)
	if err != nil {
		return nil, err
	}

	return convertToLapModels(laps)
}
