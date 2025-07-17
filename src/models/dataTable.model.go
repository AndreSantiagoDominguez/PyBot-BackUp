package models

type DataTable struct {
	Table_name string    `json:"table_name"`
	Data  interface{}    `json:"data"`
}