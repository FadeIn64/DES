package consumers

import (
	"DAS/models"
	"encoding/json"
	"log"

	"DAS/internal/repositories"
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
		var lap models.Lap
		if err := json.Unmarshal(message.Value, &lap); err != nil {
			log.Printf("Error unmarshaling message: %v\n", err)
			continue
		}

		analysis, err := h.repo.ProcessLap(session.Context(), lap)
		if err != nil || analysis == nil {
			log.Printf("Error processing lap: %v\n", err)
			session.MarkMessage(message, "")
		}

		log.Printf(
			"Lap analysis: Current=%.3fs, Avg=%.3fs(%.1f%%), SegmentAvg=%.3fs(%d laps)\n",
			analysis.CurrentLapTime,
			analysis.AverageLapTime,
			analysis.ComparisonWithAvg,
			analysis.AverageSegmentPace,
			analysis.LapsInSegment,
		)

		session.MarkMessage(message, "")
	}
	return nil
}
