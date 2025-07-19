package controllers

import (
	"PyBot-BackUp/src/connections"
	"PyBot-BackUp/src/models"
	"fmt"
	"time"
)

type PostgreSQL struct {
	conn *connections.ConnPostgreSQL
}

func NewPostgreSQL() *PostgreSQL {
	conn := connections.GetDBPool()

	if conn.Err != "" {
		fmt.Printf("Error al configurar el pool de conexiones: %v", conn.Err)
	}

	return &PostgreSQL{conn: conn}
}

func (postgre *PostgreSQL) InsertPeriod( wp models.WorkPeriod) (int, error) {
	fmt.Printf("D: %v", wp)
	query := `INSERT INTO work_periods (period_id, start_hour, end_hour, day_work, prototype_id, backup)
	      VALUES ($1,$2,$3,$4,$5, $6)
		  RETURNING period_id`

	const layout = "2006-01-02T15:04:05.999999Z07:00"
	startT, err := time.Parse(layout, wp.Start_hour)
	if err != nil {
		fmt .Printf("Error al ejecutar el parseo de la hora de inicio: %v", err)
		return 0, err

	}

	endT, err := time.Parse(layout, wp.End_hour)
	if err != nil {
		fmt .Printf("Error al ejecutar el parseo de la hora de termino: %v", err)
		return 0, err

	}

	var id int 

	err = postgre.conn.DB.QueryRow(query, wp.Period_id, startT, endT, wp.Day_work, wp.Prototype_id, true).Scan(&id)
	if err != nil {
		fmt.Printf("Error al insertar el periodo: %v", err)
		return 0, err
	}

	return id, nil
}

func (postgre *PostgreSQL) InsertReading(r models.Reading, prototype_id string)(error){
	fmt.Printf("E: %v", prototype_id)
	query := `INSERT INTO readings (period_id, distance_traveled, weight_waste, prototype_id)
			  VALUES ($1, $2, $3, $4)
			  RETURNING period_id`

	var id int 
	
	err := postgre.conn.DB.QueryRow(query,r.Period_id, r.Distance_traveled, r.Weight_waste, prototype_id).Scan(&id)
	if err != nil {
		fmt.Printf("Error al ejecutar la insercion de reading: %v", err)
		return err
	}

	return nil
		
}

func (postgres *PostgreSQL)	InsertWasteCollectionRegister(wc models.WasteCollection, prototype_id string) (int, error) {
	query := `INSERT INTO waste_collection (waste_collection_id, period_id, amount, waste_id, prototype_id)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING waste_collection_id`

	var id int		  
	
	err := postgres.conn.DB.QueryRow(query,wc.Waste_collection_id , wc.Period_id, wc.Amount, wc.Waste_id, prototype_id).Scan(&id)
	if err != nil {
		fmt.Printf("Error al ejecutar WasteCollectionRegister: %v", err)
		return 0, err
	}

	return id, nil
}

func (postgres *PostgreSQL)	InsertWeightData(w models.WeightData, prototype_id string) (int, error) {
	query := `INSERT INTO weight_data (weight_data_id, period_id, hour_period, weight, prototype_id)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING weight_data_id` 

	var id int
	
	const layout = "2006-01-02T15:04:05.999999Z07:00"
	startT, err := time.Parse(layout, w.Hour_period)
	if err != nil {
		fmt.Printf("Error al ejecutar parsear la hora de periodo: %v", err)
		return 0, err
	}

	err = postgres.conn.DB.QueryRow(query, w.Weight_data_id, w.Period_id, startT, w.Weight, prototype_id).Scan(&id)
	if err != nil {
		fmt.Printf("Error al ejecutarr la consulta WeightRegister: %v", err)
		return 0, err
	}
	
	return id, nil
}

func (postgres *PostgreSQL)	InsertGPSData(gps models.GPSData, prototype_id string) (int, error) {
	query := `INSERT INTO gps_data ( gps_data_id, period_id, latitude, longitude, altitude, speed, date_gps, hour_UTC, prototype_id)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			  RETURNING gps_data_id`

	var id int

	const layout = "2006-01-02T15:04:05.999999Z07:00"
	startT, err := time.Parse(layout, gps.Hour_UTC)
	if err != nil {
		fmt.Printf("Error al ejecutar parsear la hora UTC: %v", err)
		return 0, err
	}

	err = postgres.conn.DB.QueryRow(
		query,
		gps.Gps_data_id, 
		gps.Period_id,
		gps.Latitude,
		gps.Longitude,
		gps.Altitude,
		gps.Speed,
		gps.Date_gps,
		startT,
		prototype_id,
	).Scan(&id)	

	if err != nil {
		fmt.Printf("Error al ejecutarr GPSRegister: %v", err)
		return 0, err
	}

	return id, nil
}




