package main

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"

	"DAS/config"
	"DAS/internal/app"
)

func main() {

	cfg, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}

	application := app.NewApp(cfg)
	defer application.Close()

	consumer := setupKafkaConsumer(cfg)
	defer consumer.Close()

	prometheus.MustRegister(application.Exporter)

	go func() {
		log.Printf("Starting metrics server at :%s", application.Cfg.ServerPort)
		if err := application.Server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	go runConsumer(context.Background(), consumer, application.LapHandler, application)

	waitForShutdown()
}

func setupKafkaConsumer(cfg *config.Config) sarama.ConsumerGroup {
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRange(),
	}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup([]string{cfg.KafkaBroker}, cfg.KafkaGroupID, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	return consumer
}

func runConsumer(ctx context.Context, consumer sarama.ConsumerGroup, handler sarama.ConsumerGroupHandler, app *app.App) {
	for {
		if err := consumer.Consume(ctx, []string{app.Cfg.KafkaTopic}, handler); err != nil {
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
