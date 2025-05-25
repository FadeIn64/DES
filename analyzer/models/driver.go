package models

import (
	"time"
)

type Driver struct {
	DriverNumber int       `json:"driver_number"`
	TeamKey      int       `json:"team_key,omitempty"`
	FullName     string    `json:"full_name"`
	Abbreviation string    `json:"abbreviation"`
	Country      string    `json:"country"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	Description  string    `json:"description"`
}

type DriversRaceData struct {
	Position                 int         `json:"position"`
	MeetingKey               int         `json:"meeting_key"`
	SessionKey               int         `json:"session_key"`
	DriverNumber             int         `json:"driver_number"`
	LapNumber                int         `json:"lap_number"`
	Interval                 float64     `json:"interval"`
	PredictionLapsToOvertake int         `json:"prediction_laps_to_overtake"`
	LastLapDuration          float64     `json:"last_lap_duration"`
	Pitsops                  int         `json:"pitsops"`
	LastPitLap               interface{} `json:"last_pit_lap"`
	FullName                 string      `json:"full_name"`
	Abbreviation             string      `json:"abbreviation"`
	TeamName                 string      `json:"team_name"`
	Color                    string      `json:"color"`
}
