/*
	Hans Nicolaus
	29 Aug 2023
*/

package model

import (
	"github.com/IBM/sarama"
	"log"
)

type Consumer struct {
	Handler func(message []byte) error
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	_ = consumer
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	_ = consumer
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if message == nil {
			return nil
		}

		if err := consumer.Handler(message.Value); err != nil {
			log.Printf("[Error][Kafka] failed processing message: %v", err)
			return err
		}

		session.MarkMessage(message, "")
	}
	return nil
}
