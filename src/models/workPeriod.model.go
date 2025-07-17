package models

type WorkPeriod struct {
	Period_id    int    `json:"period_id"`
	Start_hour   string `json:"start_hour"`
	End_hour     string `json:"end_hour"`
	Day_work     string `json:"day_work"`
	Prototype_id string `json:"prototype_id"`
	BackUp       bool   `json:"backUp"`
}