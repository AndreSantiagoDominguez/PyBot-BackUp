package models

type WasteCollection struct {
	Waste_collection_id int `json:"waste_collection_id"`
	Period_id           int `json:"period_id"`
	Amount              int `json:"amount"`
	Waste_id            int `json:"waste_id"`
}