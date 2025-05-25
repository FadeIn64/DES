package repositories

import (
	"DAS/internal/repositories/db"
	"DAS/models"
	"context"
	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MeetingRepository struct {
	db      *pgxpool.Pool
	manager trm.Manager
}

func NewMeetingRepository(db *pgxpool.Pool, manager trm.Manager) *MeetingRepository {
	return &MeetingRepository{db: db, manager: manager}
}

func (r *MeetingRepository) GetAll(ctx context.Context) ([]*models.Meeting, error) {
	q := db.New(r.db)

	res, err := q.GetMeetings(ctx)
	if err != nil {
		return nil, err
	}
	meetings := make([]*models.Meeting, len(res))
	for i, meeting := range res {
		meetings[i] = &models.Meeting{
			MeetingKey:    int(meeting.MeetingKey),
			Name:          meeting.Name,
			Description:   meeting.Description,
			Circuit:       meeting.Circuit,
			Location:      meeting.Location,
			StartDate:     meeting.StartDate.Time,
			EndDate:       meeting.EndDate.Time,
			Year:          int(meeting.Year),
			DashboardLink: meeting.DashboardLink.String,
		}
	}

	return meetings, nil
}

func (r *MeetingRepository) GetMeeting(ctx context.Context, meetingKey int) (*models.Meeting, error) {
	q := db.New(r.db)

	res, err := q.GetMeetingByKey(ctx, int64(meetingKey))
	if err != nil {
		return nil, err
	}

	return &models.Meeting{
		MeetingKey:    int(res.MeetingKey),
		Name:          res.Name,
		Description:   res.Description,
		Circuit:       res.Circuit,
		Location:      res.Location,
		StartDate:     res.StartDate.Time,
		EndDate:       res.EndDate.Time,
		Year:          int(res.Year),
		DashboardLink: res.DashboardLink.String,
	}, nil
}

func (r *MeetingRepository) GetDriversStatsByMeeting(ctx context.Context, meetingKey int) ([]*models.DriversRaceData, error) {
	q := db.New(r.db)

	res, err := q.GetDriversRaceDataByMeeting(ctx, int32(meetingKey))
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
