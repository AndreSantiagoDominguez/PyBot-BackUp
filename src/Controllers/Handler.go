package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"PyBot-BackUp/src/models"
)

// Handler agrupa la conexión y operaciones a ejecutar en transacción
type Handler struct {
	postgres *PostgreSQL
}

// NewHandler crea un nuevo handler con conexión a PostgreSQL
func NewHandler() *Handler {
	postgres := NewPostgreSQL()
	return &Handler{postgres: postgres}
}

// Send inserta en la base de datos todos los registros recibidos
func (h *Handler) Send(tables []models.DataTable) (err error) {
	var prototype_id string 
	for _, table := range tables {
		if table.Data == nil {
			continue
		}
		
		switch table.Table_name {
		case "work_periods":
			var periods []models.WorkPeriod
			if err = mapToStruct(table.Data, &periods); err != nil {
				return fmt.Errorf("error mapeando work_periods: %w", err)
			}

			prototype_id = periods[0].Prototype_id

			for _, p := range periods {
				if _, err = h.postgres.InsertPeriod(p); err != nil {
					return fmt.Errorf("error insertando periodo: %w", err)
				}
				log.Printf("Periodo insertado ID: %d", p.Period_id)
			}

		case "readings":
			var readings []models.Reading
			if err = mapToStruct(table.Data, &readings); err != nil {
				return fmt.Errorf("error mapeando readings: %w", err)
			}
			for _, r := range readings {
				if err = h.postgres.InsertReading(r, prototype_id); err != nil {
					return fmt.Errorf("error insertando reading: %w", err)
				}
			}

		case "waste_collection":
			var wastes []models.WasteCollection
			if err = mapToStruct(table.Data, &wastes); err != nil {
				return fmt.Errorf("error mapeando waste_collection: %w", err)
			}
			for _, w := range wastes {
				if _, err = h.postgres.InsertWasteCollectionRegister(w, prototype_id); err != nil {
					return fmt.Errorf("error insertando waste_collection: %w", err)
				}
			}

		case "weight_data":
			var weights []models.WeightData
			if err = mapToStruct(table.Data, &weights); err != nil {
				return fmt.Errorf("error mapeando weight_data: %w", err)
			}
			for _, w := range weights {
				if _, err = h.postgres.InsertWeightData(w, prototype_id); err != nil {
					return fmt.Errorf("error insertando weight_data: %w", err)
				}
			}

		case "gps_data":
			var gpsList []models.GPSData
			if err = mapToStruct(table.Data, &gpsList); err != nil {
				return fmt.Errorf("error mapeando gps_data: %w", err)
			}
			for _, g := range gpsList {
				if _, err = h.postgres.InsertGPSData(g, prototype_id); err != nil {
					return fmt.Errorf("error insertando gps_data: %w", err)
				}
			}

		default:
			log.Printf("Tabla no reconocida: %s", table.Table_name)
		}
	}

	log.Println("Todas las tablas procesadas correctamente")
	return
}

// mapToStruct convierte un interface{} en target usando JSON
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
