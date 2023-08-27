package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"bibit.id/challenge/handler"
	"github.com/IBM/sarama"
)

func serveKafka(handler *handler.Handler) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	brokers := []string{"localhost:9092"}

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("[Error][Kafka] Failed creating consumer: %v", err)
	}
	defer func() {
		if err = consumer.Close(); err != nil {
			log.Printf("[Error][Kafka] Failed closing consumer: %v", err)
		}
	}()

	topic := "bibit_challenge_1"

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Printf("[Error][Kafka] Failed creating partition consumer: %v", err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Printf("[Error][Kafka] Failed closing partition consumer: %v", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	log.Print("[Kafka] Listening on port 9092")

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			err := handler.ProcessStockTransaction(context.Background(), msg.Value)
			if err != nil {
				log.Printf("[Error][Kafka] Failed ProcessStockTransaction: %s", err.Error())
			}
		case err := <-partitionConsumer.Errors():
			log.Printf("[Error][Kafka] Failed consuming by partition consumer: %v", err)
		case <-signals:
			return
		}
	}
}
