package receivers

import (
	"DES/extractors"
	"DES/models"
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"log"
	"time"
)

type Receiver interface {
	ReceiveData(ctx context.Context, start time.Time) chan error
	Close() error
}

type lapReceiver struct {
	producer  sarama.SyncProducer
	topic     string
	extractor extractors.LapExtractor
	canceled  bool
}

func NewLapReceiver(producer sarama.SyncProducer, topic string, extractor extractors.LapExtractor) Receiver {
	return &lapReceiver{
		producer:  producer,
		topic:     topic,
		extractor: extractor,
		canceled:  false,
	}
}

func (l *lapReceiver) ReceiveData(ctx context.Context, start time.Time) chan error {
	var errChan = make(chan error)

	go func(extractor extractors.LapExtractor) {

		lastNow := time.Now()
		fromRaceTime := start

		for !l.canceled {
			curNow := time.Now()
			diff := curNow.Sub(lastNow)

			toRaceTime := fromRaceTime.Add(diff)

			laps, err := extractor.ExtractLaps(ctx, fromRaceTime, toRaceTime)
			if err != nil {
				errChan <- err
				return
			}
			log.Printf("Laps: %v", laps)

			for _, lap := range laps {
				err = l.sendData(&lap)
				if err != nil {
					errChan <- err
					return
				}
			}

			lastNow = curNow
			fromRaceTime = toRaceTime
			<-time.After(2 * time.Second)

		}
	}(l.extractor)

	return errChan
}

func (l *lapReceiver) Close() error {
	l.canceled = true
	return nil
}

func (l *lapReceiver) sendData(lap *models.Lap) error {
	bytes, err := json.Marshal(lap)
	if err != nil {
		return err
	}

	message := &sarama.ProducerMessage{
		Topic: l.topic,
		Key:   sarama.StringEncoder(uuid.New().String()),
		Value: sarama.ByteEncoder(bytes),
	}

	_, _, err = l.producer.SendMessage(message)
	return err
}
