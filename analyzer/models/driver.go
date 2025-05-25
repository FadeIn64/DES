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
	Position                 int
	MeetingKey               int
	SessionKey               int
	DriverNumber             int
	LapNumber                int
	Interval                 float64
	PredictionLapsToOvertake int
	LastLapDuration          float64
	Pitsops                  int
	LastPitLap               interface{}
	FullName                 string
	Abbreviation             string
	Name                     string
	Color                    string
}
