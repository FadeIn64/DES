package consumers

import (
	"DAS/internal/repositories"
	"DAS/models"
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
)

type DriverHandler struct {
	repo *repositories.DriverRepository
}

func NewDriverHandler(repo *repositories.DriverRepository) *DriverHandler {
	return &DriverHandler{repo: repo}
}

func (d *DriverHandler) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Kafka drivers consumer setup completed")
	return nil
}

func (d *DriverHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Kafka drivers consumer cleanup completed")
	return nil
}

func (d *DriverHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var driver models.Driver
		if err := json.Unmarshal(message.Value, &driver); err != nil {
			log.Printf("Error on unmarshalling message: %s\n", err)
			continue
		}

		if err := d.repo.SaveDriver(session.Context(), &driver); err != nil {
			log.Printf("Error saving driver: %s\n", err)
			continue
		}

		session.MarkMessage(message, "")
	}
	return nil
}
