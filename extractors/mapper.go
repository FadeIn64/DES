package extractors

import (
	"DES/db"
	"DES/models"
)

func convertToLapModel(entity db.Lap) (*models.Lap, error) {
	return &models.Lap{
		MeetingKey:      entity.MeetingKey,
		SessionKey:      entity.SessionKey,
		DriverNumber:    entity.DriverNumber.Int32,
		I1Speed:         entity.I1Speed.Int32,
		I2Speed:         entity.I2Speed.Int32,
		StSpeed:         entity.StSpeed.Int32,
		DateStart:       entity.DateStart.Time,
		LapDuration:     entity.LapDuration.Float64,
		IsPitOutLap:     entity.IsPitOutLap.Bool,
		DurationSector1: entity.DurationSector1.Float64,
		DurationSector2: entity.DurationSector2.Float64,
		DurationSector3: entity.DurationSector3.Float64,
		SegmentsSector1: entity.SegmentsSector1,
		SegmentsSector2: entity.SegmentsSector2,
		SegmentsSector3: entity.SegmentsSector3,
		LapNumber:       entity.LapNumber,
	}, nil
}

func convertToLapModels(entity []db.Lap) ([]models.Lap, error) {
	result := make([]models.Lap, len(entity))
	for i, lap := range entity {
		lapModel, err := convertToLapModel(lap)
		if err != nil {
			return nil, err
		}
		result[i] = *lapModel
	}
	return result, nil
}
