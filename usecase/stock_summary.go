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

	summary, err := uc.stockRepo.GetStockSummary(ctx, model.GetStockSummaryRequest{
		StockCode: stockCode,
		Date:      transactionDate,
	})
	if err != nil {
		return err
	}

	if summary == (model.Summary{}) {
		summary.StockCode = stockCode
	}

	if isUpdated, updatedSummary := summary.ApplyTransaction(transaction); isUpdated {
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
