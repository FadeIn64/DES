// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	ExecCompletedLapByDriver(ctx context.Context, arg ExecCompletedLapByDriverParams) (int64, error)
	GetAllTeams(ctx context.Context) ([]Team, error)
	GetAverageLapTime(ctx context.Context, arg GetAverageLapTimeParams) (float64, error)
	GetCurrentMeeting(ctx context.Context) (Meeting, error)
	GetCurrentSegmentPace(ctx context.Context, arg GetCurrentSegmentPaceParams) (GetCurrentSegmentPaceRow, error)
	GetDriverByNumber(ctx context.Context, driverNumber int64) (Driver, error)
	GetDriverByPosition(ctx context.Context, arg GetDriverByPositionParams) (DriversStatsWithPosition, error)
	GetDriverStats(ctx context.Context, arg GetDriverStatsParams) (DriversStatsWithPosition, error)
	GetDriversByTeam(ctx context.Context, teamKey pgtype.Int4) ([]Driver, error)
	GetDriversRaceDataByDriver(ctx context.Context, driverNumber int32) ([]FullDriverDatum, error)
	GetDriversRaceDataByMeeting(ctx context.Context, meetingKey int32) ([]FullDriverDatum, error)
	GetDriversStats(ctx context.Context, arg GetDriversStatsParams) ([]DriversStatsWithPosition, error)
	GetLap(ctx context.Context, arg GetLapParams) (Lap, error)
	GetMeetingByKey(ctx context.Context, meetingKey int64) (Meeting, error)
	GetMeetings(ctx context.Context) ([]Meeting, error)
	GetTeamByID(ctx context.Context, teamKey int64) (Team, error)
	MoveCompleteLap(ctx context.Context, arg MoveCompleteLapParams) error
	UpsertDrivers(ctx context.Context, arg UpsertDriversParams) error
	UpsertDriversInterval(ctx context.Context, arg UpsertDriversIntervalParams) error
	UpsertLap(ctx context.Context, arg UpsertLapParams) error
	UpsertMeeting(ctx context.Context, arg UpsertMeetingParams) error
	UpsertTeam(ctx context.Context, arg UpsertTeamParams) error
}

var _ Querier = (*Queries)(nil)
