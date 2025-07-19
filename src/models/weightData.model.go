package models

type WeightData struct {
	Weight_data_id int     `json:"weight_data_id"`
	Period_id      int     `json:"period_id"`
	Hour_period    string  `json:"hour_period"`
	Weight         float32 `json:"weight"`
	Prototype_id   string  `json:"prototype_id"`
}