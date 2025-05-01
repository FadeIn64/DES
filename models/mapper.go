package models

import "time"

func F1LapToSectorLapsWithTiming(f1Lap F1Lap) []Lap {
	return []Lap{
		// 1-я структура: только Sector 1
		{
			MeetingKey:     f1Lap.MeetingKey,
			SessionKey:     f1Lap.SessionKey,
			DriverNumber:   f1Lap.DriverNumber,
			DateStart:      f1Lap.DateStart,
			LapDuration:    0,
			LapNumber:      f1Lap.LapNumber,
			SectorDuration: []float64{f1Lap.DurationSector1, 0, 0},
			InfoTime:       f1Lap.DateStart.Add(time.Duration(f1Lap.DurationSector1 * float64(time.Second))),
		},
		// 2-я структура: Sector 1 + Sector 2
		{
			MeetingKey:     f1Lap.MeetingKey,
			SessionKey:     f1Lap.SessionKey,
			DriverNumber:   f1Lap.DriverNumber,
			DateStart:      f1Lap.DateStart,
			LapDuration:    0,
			LapNumber:      f1Lap.LapNumber,
			SectorDuration: []float64{f1Lap.DurationSector1, f1Lap.DurationSector2, 0},
			InfoTime:       f1Lap.DateStart.Add(time.Duration((f1Lap.DurationSector1 + f1Lap.DurationSector2) * float64(time.Second))),
		},
		// 3-я структура: все три сектора
		{
			MeetingKey:     f1Lap.MeetingKey,
			SessionKey:     f1Lap.SessionKey,
			DriverNumber:   f1Lap.DriverNumber,
			DateStart:      f1Lap.DateStart,
			LapDuration:    f1Lap.LapDuration,
			LapNumber:      f1Lap.LapNumber,
			SectorDuration: []float64{f1Lap.DurationSector1, f1Lap.DurationSector2, f1Lap.DurationSector3},
			InfoTime:       f1Lap.DateStart.Add(time.Duration((f1Lap.DurationSector1 + f1Lap.DurationSector2 + f1Lap.DurationSector3) * float64(time.Second))),
		},
	}
}
