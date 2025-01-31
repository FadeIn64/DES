package models

import "time"

type Lap struct {
	MeetingKey      int32     `json:"meeting_key"`
	SessionKey      int32     `json:"session_key"`
	DriverNumber    int32     `json:"driver_number"`
	I1Speed         int32     `json:"i1_speed"`
	I2Speed         int32     `json:"i2_speed"`
	StSpeed         int32     `json:"st_speed"`
	DateStart       time.Time `json:"date_start"`
	LapDuration     float64   `json:"lap_duration"`
	IsPitOutLap     bool      `json:"is_pit_out_lap"`
	DurationSector1 float64   `json:"duration_sector_1"`
	DurationSector2 float64   `json:"duration_sector_2"`
	DurationSector3 float64   `json:"duration_sector_3"`
	SegmentsSector1 []int32   `json:"segments_sector_1"`
	SegmentsSector2 []int32   `json:"segments_sector_2"`
	SegmentsSector3 []int32   `json:"segments_sector_3"`
	LapNumber       int32     `json:"lap_number"`
}
