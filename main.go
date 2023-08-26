package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"

	"bibit.id/challenge/model"
	"bibit.id/challenge/repo"
	"bibit.id/challenge/usecase"
)

func main() {
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

	stockRepo := repo.New()
	uc := usecase.New(stockRepo)

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			handler(uc, msg.Value)
		case err := <-partitionConsumer.Errors():
			log.Printf("Error: %v", err)
		case <-signals:
			return
		}
	}
}

func handler(usecase *usecase.Usecase, data []byte) {
	input := model.Input{}
	if err := json.Unmarshal(data, &input); err != nil {
		fmt.Println("UNMARSHAL ERROR")
		return
	}

	transaction, err := input.ToTransaction()
	if err != nil {
		fmt.Println("Error converting to ChangeRecord:", err)
		return
	}

	err = usecase.UpdateStockSummary(transaction)
	if err != nil {
		fmt.Println(err.Error())
	}
}
