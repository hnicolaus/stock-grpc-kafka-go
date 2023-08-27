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
	// Get stock summary by stockCode and date if already exists
	summaryResult, err := uc.stockRepo.GetStockSummary(ctx, model.GetStockSummaryRequest{
		StockCode: transaction.StockCode,
		FromDate:  transaction.Date,
		ToDate:    transaction.Date,
	})
	if err != nil {
		return err
	}

	summary := model.Summary{}
	if len(summaryResult) > 0 {
		summary = summaryResult[0]
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

func (uc *Usecase) GetStockSummary(ctx context.Context, request model.GetStockSummaryRequest) ([]model.Summary, error) {
	return uc.stockRepo.GetStockSummary(ctx, request)
}
