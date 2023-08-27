/*
	Hans Nicolaus
	26 Aug 2023
*/

package usecase

import (
	"context"

	"bibit.id/challenge/model"
)

func (uc *Usecase) UpdateStockSummary(ctx context.Context, transaction model.Transaction) error {
	var (
		stockCode       = transaction.StockCode
		transactionDate = transaction.Date
	)

	// Get stock summary by stockCode and date if alrady exists
	summary, err := uc.stockRepo.GetStockSummary(ctx, model.GetStockSummaryRequest{
		StockCode: stockCode,
		Date:      transactionDate,
	})
	if err != nil {
		return err
	}

	// Update stock summary data based on the transaction
	isUpdated, updatedSummary := summary.ApplyTransaction(transaction)

	if isUpdated {
		// Persist updated stock summary to our data store
		err = uc.stockRepo.SetStockSummary(ctx, updatedSummary)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *Usecase) GetStockSummary(ctx context.Context, request model.GetStockSummaryRequest) (model.Summary, error) {
	return uc.stockRepo.GetStockSummary(ctx, request)
}
