package main

import (
	"DES/extractors"
	"DES/receivers"
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"log"
	"os/signal"
	"syscall"
	"time"
)

const (
	dbUrl      = "postgres://username:password@localhost:5432/des"
	dateString = "2024-07-21 13:05:45.000000"
	layout     = "2006-01-02 15:04:05.999999"
	kafkaUrl   = "localhost:9092"
	lapTopic   = "laps"
)

func main() {

	mainCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	//init db
	dbConf, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.NewWithConfig(mainCtx, dbConf)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	db := stdlib.OpenDBFromPool(pool)

	//goose
	err = goose.SetDialect("postgres")
	if err != nil {
		log.Fatal(err)
	}

	err = goose.Up(db, "resources/changelog")
	if err != nil {
		log.Fatal(err)
	}

	//producer
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{kafkaUrl}, config)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	startDate, err := getStartDate()
	if err != nil {
		log.Fatal(err)
	}

	//extractor
	lapExtractor := extractors.NewLapExtractor(pool)

	//receivers
	lapReceiver := receivers.NewLapReceiver(producer, lapTopic, lapExtractor)

	lapReceiverErrChan := lapReceiver.ReceiveData(mainCtx, startDate)
	defer lapReceiver.Close()

	select {
	case <-mainCtx.Done():
		return
	case err := <-lapReceiverErrChan:
		log.Fatal(err)
	}
}

func getStartDate() (time.Time, error) {

	// Парсим строку в time.Time
	t, err := time.Parse(layout, dateString)
	if err != nil {
		fmt.Println("Ошибка при парсинге:", err)
		return time.Time{}, err
	}
	return t, nil
}
