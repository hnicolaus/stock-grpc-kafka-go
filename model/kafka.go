/*
	Hans Nicolaus
	29 Aug 2023
*/

package model

import (
	"github.com/IBM/sarama"
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

		_ = consumer.Handler(message.Value)

		session.MarkMessage(message, "")
	}
	return nil
}
