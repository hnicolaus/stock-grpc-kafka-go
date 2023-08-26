/*
	Hans Nicolaus
	26 Aug 2023
*/

package usecase

import (
	"bibit.id/challenge/model"
)

func (uc *Usecase) UpdateStockSummary(transaction model.Transaction) error {
	var (
		stockCode       = transaction.StockCode
		transactionDate = transaction.Date
	)

	summary, err := uc.stockRepo.GetStockSummary(stockCode, transactionDate)
	if err != nil {
		return err
	}

	if summary == (model.Summary{}) {
		summary.StockCode = stockCode
	}

	if isUpdated, updatedSummary := summary.ApplyTransaction(transaction); isUpdated {
		err = uc.stockRepo.SetStockSummary(stockCode, transactionDate, updatedSummary)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *Usecase) GetStockSummary(transaction model.Transaction) (model.Summary, error) {
	var (
		stockCode       = transaction.StockCode
		transactionDate = transaction.Date
	)

	return uc.stockRepo.GetStockSummary(stockCode, transactionDate)
}
