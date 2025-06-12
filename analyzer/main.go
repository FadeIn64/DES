package main

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"

	"DAS/config"
	"DAS/internal/app"
)

func main() {

	mainCtx, mainCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer mainCancel()

	cfg, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}

	application := app.NewApp(cfg)
	defer application.Close()

	prometheus.MustRegister(application.Exporter)

	go func() {
		log.Printf("Starting metrics server at :%s", application.Cfg.ServerPort)
		if err := application.Server.ListenAndServe(); err != nil {
			mainCancel()
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	runConsumers(mainCtx, application)

	<-mainCtx.Done()

}

func setupKafkaConsumer(cfg *config.Config) sarama.ConsumerGroup {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Version = sarama.V2_5_0_0
	kafkaConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRange(),
	}
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup([]string{cfg.KafkaBroker}, cfg.KafkaGroupID, kafkaConfig)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	return consumer
}

func runConsumers(ctx context.Context, app *app.App) {
	go runConsumer(ctx, app, app.LapHandler, app.Cfg.KafkaTopic, "laps")
	go runConsumer(ctx, app, app.MeetingHandler, "meetings", "meetings")
	go runConsumer(ctx, app, app.DriverHandler, "drivers", "drivers")
	go runConsumer(ctx, app, app.TeamHandler, "teams", "teams")
}

func runConsumer(ctx context.Context, app *app.App, handler sarama.ConsumerGroupHandler, topic, name string) {
	consumer := setupKafkaConsumer(app.Cfg)
	for {
		if err := consumer.Consume(ctx, []string{topic}, handler); err != nil {
			log.Printf("Error from %s consumer: %v", name, err)
		}
		if ctx.Err() != nil {
			return
		}
	}
}
