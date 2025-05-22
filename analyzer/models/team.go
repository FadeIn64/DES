package models

type Team struct {
	TeamKey     int    `json:"team_key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Country     string `json:"country"`
}
