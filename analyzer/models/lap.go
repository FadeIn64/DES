package models

import "time"

type Lap struct {
	MeetingKey     int32     `json:"meeting_key"`
	SessionKey     int32     `json:"session_key"`
	DriverNumber   int32     `json:"driver_number"`
	DateStart      time.Time `json:"date_start"`
	LapDuration    float64   `json:"lap_duration"`
	LapNumber      int32     `json:"lap_number"`
	SectorDuration []float64 `json:"sector_duration"`
	InfoTime       time.Time `json:"info_time"`
	IsPitOutLap    bool      `json:"is_pit_out_lap"`
}

type LapAnalysis struct {
	MeetingKey         int32
	SessionKey         int32
	DriverNumber       int32
	LapNumber          int32
	CurrentLapTime     float64
	AverageLapTime     float64
	AverageSegmentPace float64
	LapsInSegment      int
	ComparisonWithAvg  float64 // Разница текущего круга со средним (%)
	PositionTrend      string  // Тренд позиции ("improving", "stable", "declining")
}
