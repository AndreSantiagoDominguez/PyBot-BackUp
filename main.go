package main

import (
	controllers "PyBot-BackUp/src/Controllers"
	"PyBot-BackUp/src/connections"
	"PyBot-BackUp/src/models"
	"encoding/json"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	
	rabbitMQ := connections.NewRabbitMQ()
	handler := controllers.NewHandler()

	msgs := rabbitMQ.GetMessages()

	for d := range msgs {
		var tables []models.DataTable
		if err := json.Unmarshal(d.Body, &tables); err != nil {
			log.Printf("Error deserializando: %v | Mensaje: %s", err, string(d.Body))
			continue
		}
		
		log.Printf("Recibido %d tablas para procesar", len(tables))
		
		if err := handler.Send(tables); err != nil {
			log.Printf("Error procesando mensaje: %v", err)
		} else {
			d.Ack(false) 
		}
	}
}