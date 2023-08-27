package handler

import (
	"context"
	"encoding/json"

	"bibit.id/challenge/model"
)

func (h *Handler) ProcessStockTransaction(ctx context.Context, data []byte) error {
	input := model.Input{}
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}

	transaction, err := input.ToTransaction()
	if err != nil {
		return err
	}

	return h.stockUsecase.UpdateStockSummary(ctx, transaction)
}
