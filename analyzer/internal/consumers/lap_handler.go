package consumers

import (
	"encoding/json"
	"log"

	"DAS/internal/repositories"
	"DAS/internal/repositories/db"
	"github.com/IBM/sarama"
)

type LapHandler struct {
	repo *repositories.LapRepository
}

func NewLapHandler(repo *repositories.LapRepository) *LapHandler {
	return &LapHandler{repo: repo}
}

func (h *LapHandler) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Kafka consumer setup completed")
	return nil
}

func (h *LapHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Kafka consumer cleanup completed")
	return nil
}

func (h *LapHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var lap db.Lap
		if err := json.Unmarshal(message.Value, &lap); err != nil {
			log.Printf("Error unmarshaling message: %v\n", err)
			continue
		}

		if err := h.repo.ProcessLap(session.Context(), lap); err != nil {
			log.Printf("Error processing lap: %v\n", err)
		}

		session.MarkMessage(message, "")
	}
	return nil
}
