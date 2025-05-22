package models

import "time"

type Meeting struct {
	MeetingKey    int       `json:"meeting_key"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Circuit       string    `json:"circuit"`
	Location      string    `json:"location"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	Year          int       `json:"year"`
	DashboardLink string    `json:"dashboard_link"`
}
