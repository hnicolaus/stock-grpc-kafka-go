package main

import (
	"context"
	"fmt"
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
		log.Fatalf("Error creating consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("Error closing consumer: %v", err)
		}
	}()

	topic := "bibit_challenge_1"

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %v", err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Printf("Error closing partition consumer: %v", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	fmt.Println("Kafka server: listening on port 9092")
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			handler.ProcessStockTransaction(context.TODO(), msg.Value)
		case err := <-partitionConsumer.Errors():
			log.Printf("Error: %v", err)
		case <-signals:
			return
		}
	}
}
