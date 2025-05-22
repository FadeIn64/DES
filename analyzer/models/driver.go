package models

import "time"

type Driver struct {
	DriverNumber int       `json:"driver_number"`
	TeamKey      int       `json:"team_key,omitempty"`
	FullName     string    `json:"full_name"`
	Abbreviation string    `json:"abbreviation"`
	Country      string    `json:"country"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	Description  string    `json:"description"`
}
