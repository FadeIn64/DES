package consumers

import (
	"DAS/internal/repositories"
	"DAS/models"
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
)

type TeamHandler struct {
	repo *repositories.TeamRepository
}

func NewTeamHandler(repo *repositories.TeamRepository) *TeamHandler {
	return &TeamHandler{repo: repo}
}

func (t *TeamHandler) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Kafka teams consumer setup completed")
	return nil
}

func (t *TeamHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Kafka teams consumer cleanup completed")
	return nil
}

func (t *TeamHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var team models.Team
		if err := json.Unmarshal(message.Value, &team); err != nil {
			log.Printf("Error on unmarshalling message: %s\n", err)
			continue
		}

		if err := t.repo.Save(session.Context(), &team); err != nil {
			log.Printf("Error saving team: %s\n", err)
			continue
		}

		session.MarkMessage(message, "")
	}
	return nil
}
