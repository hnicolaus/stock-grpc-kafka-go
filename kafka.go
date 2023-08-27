package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"bibit.id/challenge/handler"
	"bibit.id/challenge/model"
	"github.com/IBM/sarama"
)

func serveKafka(handler *handler.Handler) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Version = sarama.DefaultVersion

	brokers := []string{"localhost:9092"}
	groupID := "bibit_consumer_group"

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("[Error][Kafka] Failed creating consumer group: %v", err)
	}
	defer func() {
		if err = consumerGroup.Close(); err != nil {
			log.Printf("[Error][Kafka] Failed closing consumer group: %v", err)
		}
	}()

	topics := []string{"bibit_challenge_1"}

	consumer := &model.Consumer{Handler: handler.ProcessStockTransaction}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	log.Print("[Kafka] Serving on port 9092")
	for {
		if err := consumerGroup.Consume(context.Background(), topics, consumer); err != nil {
			log.Printf("[Error][Kafka] Error from consumer: %v", err)
		}

		select {
		case <-signals:
			return
		default:
		}
	}
}
