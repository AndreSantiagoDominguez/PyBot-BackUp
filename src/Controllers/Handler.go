package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"PyBot-BackUp/src/models"
)

type Handler struct {
	postgres *PostgreSQL
}

func NewHandler() *Handler {
	postgres := NewPostgreSQL()
    return &Handler{postgres: postgres}
}

func (h *Handler) Send(tables []models.DataTable) error {
	tx, err := h.postgres.conn.DB.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Transacción revertida por panic: %v", r)
		}
	}()

	// Primera pasada: procesar work_periods para obtener IDs
	periodIDs := make(map[int]bool)
	for _, table := range tables {
		if table.Table_name == "work_periods" && table.Data != nil {
			var periods []models.WorkPeriod
			if err := mapToStruct(table.Data, &periods); err != nil {
				return fmt.Errorf("error mapeando work_periods: %w", err)
			}

			for _, period := range periods {
				periodIDs[period.Period_id] = true
			}
		}
	}

	// Segunda pasada: procesar todas las tablas
	for _, table := range tables {
		if table.Data == nil {
			continue
		}

		switch table.Table_name {
		case "work_periods":
			var periods []models.WorkPeriod
			if err := mapToStruct(table.Data, &periods); err != nil {
				tx.Rollback()
				return fmt.Errorf("error mapeando work_periods: %w", err)
			}

			for _, period := range periods {


				if !periodIDs[period.Period_id] {
					insertedID, err := h.postgres.InsertPeriod(period)
					if err != nil {
						tx.Rollback()
						return fmt.Errorf("error insertando periodo: %w", err)
					}
					log.Printf("Periodo insertado ID: %d", insertedID)
				} else {
					log.Printf("Periodo ya existe ID: %d", period.Period_id)
				}
			}

		case "readings":
			var readings []models.Reading
			if err := mapToStruct(table.Data, &readings); err != nil {
				tx.Rollback()
				return fmt.Errorf("error mapeando readings: %w", err)
			}

			for _, reading := range readings {
				if !periodIDs[reading.Period_id] {
					tx.Rollback()
					return fmt.Errorf("period_id %d no encontrado para reading", reading.Period_id)
				}

				if err := h.postgres.InsertReading(reading); err != nil {
					tx.Rollback()
					return fmt.Errorf("error insertando reading: %w", err)
				}
			}

		case "waste_collection":
			var wastes []models.WasteCollection
			if err := mapToStruct(table.Data, &wastes); err != nil {
				tx.Rollback()
				return fmt.Errorf("error mapeando waste_collection: %w", err)
			}

			for _, waste := range wastes {
				if !periodIDs[waste.Period_id] {
					tx.Rollback()
					return fmt.Errorf("period_id %d no encontrado para waste collection", waste.Period_id)
				}

				if _, err := h.postgres.InsertWasteCollectionRegister(waste); err != nil {
					tx.Rollback()
					return fmt.Errorf("error insertando waste_collection: %w", err)
				}
			}

		case "weight_data":
			var weights []models.WeightData
			if err := mapToStruct(table.Data, &weights); err != nil {
				tx.Rollback()
				return fmt.Errorf("error mapeando weight_data: %w", err)
			}

			for _, weight := range weights {
				if !periodIDs[weight.Period_id] {
					tx.Rollback()
					return fmt.Errorf("period_id %d no encontrado para weight data", weight.Period_id)
				}

				if _, err := h.postgres.InsertWeightData(weight); err != nil {
					tx.Rollback()
					return fmt.Errorf("error insertando weight_data: %w", err)
				}
			}

		case "gps_data":
			var gpsData []models.GPSData
			if err := mapToStruct(table.Data, &gpsData); err != nil {
				tx.Rollback()
				return fmt.Errorf("error mapeando gps_data: %w", err)
			}

			for _, gps := range gpsData {
				if !periodIDs[gps.Period_id] {
					tx.Rollback()
					return fmt.Errorf("period_id %d no encontrado para gps data", gps.Period_id)
				}

				if _, err := h.postgres.InsertGPSData(gps); err != nil {
					tx.Rollback()
					return fmt.Errorf("error insertando gps_data: %w", err)
				}
			}

		default:
			log.Printf("Tabla no reconocida: %s", table.Table_name)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error en commit: %w", err)
	}

	log.Println("Todas las tablas procesadas correctamente")
	return nil
}

// Función auxiliar para convertir interface{} a slice de estructuras
func mapToStruct(data interface{}, target interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error serializando datos: %w", err)
	}
	
	if err := json.Unmarshal(jsonData, target); err != nil {
		return fmt.Errorf("error deserializando a estructura: %w", err)
	}
	return nil
}