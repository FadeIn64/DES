package consumers

import (
	"DAS/internal/repositories"
	"DAS/models"
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
)

type MeetingHandler struct {
	repo *repositories.MeetingRepository
}

func NewMeetingHandler(repo *repositories.MeetingRepository) *MeetingHandler {
	return &MeetingHandler{repo: repo}
}

func (m *MeetingHandler) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Kafka meetings consumer setup completed")
	return nil
}

func (m *MeetingHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Kafka meetings consumer cleanup completed")
	return nil
}

func (m *MeetingHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var meeting models.Meeting
		if err := json.Unmarshal(message.Value, &meeting); err != nil {
			log.Printf("Error on unmarshalling message: %s\n", err)
			continue
		}

		if err := m.repo.SaveMeeting(session.Context(), &meeting); err != nil {
			log.Printf("Error saving meeting: %s\n", err)
			continue
		}

		session.MarkMessage(message, "")
	}
	return nil
}
