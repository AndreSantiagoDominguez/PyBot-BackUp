package models

type GPSData struct {
	Gps_data_id int     `json:"gps_data_id"`
	Period_id   int     `json:"period_id"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Altitude    float64 `json:"altitude"`
	Speed       float64 `json:"speed"`
	Date_gps    string  `json:"date_gps"`
	Hour_UTC    string  `json:"hour_UTC"`
}