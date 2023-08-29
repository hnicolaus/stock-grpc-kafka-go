/*
	Hans Nicolaus
	29 Aug 2023
*/

package handler

import (
	"context"
	"encoding/json"
	"log"

	"bibit.id/challenge/model"
)

func (h *Handler) ProcessStockTransaction(data []byte) error {
	input := model.KafkaTransaction{}
	if err := json.Unmarshal(data, &input); err != nil {
		log.Printf("[Error][ProcessStockTransaction] error unmarshaling Transaction event: %v", err)
		return err
	}

	transaction, err := input.ToTransaction()
	if err != nil {
		log.Printf("[Error][ProcessStockTransaction] error converting Transaction data: %v", err)
		return err
	}

	err = h.stockUsecase.UpdateStockSummary(context.Background(), transaction)
	if err != nil {
		log.Printf("[Error][UpdateStockSummary] %v", err)
		return err
	}

	return nil
}
