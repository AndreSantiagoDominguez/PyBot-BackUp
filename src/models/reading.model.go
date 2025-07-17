package models

type Reading struct {
	Period_id         int     `json:"period_id"`
	Distance_traveled float32 `json:"distance_traveled"`
	Weight_waste      float32 `json:"weight_waste"`
}