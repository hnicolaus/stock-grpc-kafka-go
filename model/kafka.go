package model

import (
	"github.com/IBM/sarama"
)

type Consumer struct {
	Handler func(message []byte) error
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if message == nil {
			return nil
		}

		h.Handler(message.Value)

		session.MarkMessage(message, "")
	}
	return nil
}
