package extractors

import (
	"DES/db"
	"DES/models"
)

func convertToLapModel(entity db.Lap) (*models.Lap, error) {
	return &models.Lap{
		MeetingKey:     entity.MeetingKey,
		SessionKey:     entity.SessionKey,
		DriverNumber:   entity.DriverNumber,
		DateStart:      entity.DateStart.Time,
		LapDuration:    entity.LapDuration,
		LapNumber:      entity.LapNumber,
		SectorDuration: entity.SectorDuration,
		InfoTime:       entity.InfoTime.Time,
		IsPitOutLap:    entity.IsPitOutLap,
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
