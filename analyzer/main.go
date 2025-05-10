package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"

	"DAS/config"
	"DAS/internal/app"
)

func main() {
	cfg := &config.Config{
		KafkaBrokers: []string{"localhost:9092"},
		KafkaTopic:   "laps",
		KafkaGroupID: "lap-aggregator-group",
		PGConnString: "postgres://username:password@localhost/das?sslmode=disable",
		SectorsCount: 3,
	}

	application := app.NewApp(cfg)
	defer application.Close()

	consumer := setupKafkaConsumer(cfg)
	defer consumer.Close()

	go runConsumer(context.Background(), consumer, application.LapHandler)

	waitForShutdown()
}

func setupKafkaConsumer(cfg *config.Config) sarama.ConsumerGroup {
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRange(),
	}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(cfg.KafkaBrokers, cfg.KafkaGroupID, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	return consumer
}

func runConsumer(ctx context.Context, consumer sarama.ConsumerGroup, handler sarama.ConsumerGroupHandler) {
	for {
		if err := consumer.Consume(ctx, []string{"race-laps"}, handler); err != nil {
			log.Printf("Error from consumer: %v", err)
		}
		if ctx.Err() != nil {
			return
		}
	}
}

func waitForShutdown() {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
	log.Println("Received termination signal, initiating shutdown... ", sigterm)
}
