package repositories

import (
	"DAS/internal/repositories/db"
	"DAS/models"
	"context"
	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DriverRepository struct {
	db      *pgxpool.Pool
	manager trm.Manager
}

func NewDriverRepository(db *pgxpool.Pool, manager trm.Manager) *DriverRepository {
	return &DriverRepository{
		db:      db,
		manager: manager,
	}
}

func (repo *DriverRepository) GetDriverByNumber(ctx context.Context, number int) (*models.Driver, error) {
	q := db.New(repo.db)

	res, err := q.GetDriverByNumber(ctx, int64(number))
	if err != nil {
		return nil, err
	}
	return &models.Driver{
		DriverNumber: int(res.DriverNumber),
		TeamKey:      int(res.TeamKey.Int32),
		FullName:     res.FullName,
		Abbreviation: res.Abbreviation,
		Country:      res.Country.String,
		DateOfBirth:  res.DateOfBirth.Time,
		Description:  res.Description,
	}, nil
}

func (repo *DriverRepository) GetDriversByTeam(ctx context.Context, teamKey int) ([]*models.Driver, error) {

	q := db.New(repo.db)

	res, err := q.GetDriversByTeam(ctx, pgtype.Int4{Int32: int32(teamKey), Valid: true})
	if err != nil {
		return nil, err
	}

	drivers := make([]*models.Driver, len(res))
	for i, d := range res {
		drivers[i] = &models.Driver{
			DriverNumber: int(d.DriverNumber),
			TeamKey:      int(d.TeamKey.Int32),
			FullName:     d.FullName,
			Abbreviation: d.Abbreviation,
			Country:      d.Country.String,
			DateOfBirth:  d.DateOfBirth.Time,
			Description:  d.Description,
		}
	}

	return drivers, nil
}

func (repo *DriverRepository) GetDriversStatsByDriver(ctx context.Context, driver int) ([]*models.DriversRaceData, error) {
	q := db.New(repo.db)

	res, err := q.GetDriversRaceDataByDriver(ctx, int32(driver))
	if err != nil {
		return nil, err
	}
	meetings := make([]*models.DriversRaceData, len(res))
	for i, meeting := range res {
		meetings[i] = &models.DriversRaceData{
			Position:                 int(meeting.Position),
			MeetingKey:               int(meeting.MeetingKey),
			SessionKey:               int(meeting.SessionKey),
			DriverNumber:             int(meeting.DriverNumber),
			LapNumber:                int(meeting.LapNumber),
			Interval:                 meeting.Interval,
			PredictionLapsToOvertake: int(meeting.PredictionLapsToOvertake.Int32),
			LastLapDuration:          meeting.LastLapDuration,
			Pitsops:                  int(meeting.Pitsops.Int64),
			LastPitLap:               meeting.LastPitLap,
			FullName:                 meeting.FullName,
			Abbreviation:             meeting.Abbreviation,
			TeamName:                 meeting.Name,
			Color:                    meeting.Color.String,
		}
	}
	return meetings, nil
}
