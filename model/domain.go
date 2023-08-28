/*
	Hans Nicolaus
	29 Aug 2023
*/

package model

import (
	"time"
)

type TransactionType string

const (
	TransactionTypeA         TransactionType = "A"
	TransactionTypeP         TransactionType = "P"
	TransactionTypeE         TransactionType = "E"
	TransactionTypeUndefined TransactionType = ""
)

type Transaction struct {
	Price     int64
	Quantity  int64
	StockCode string
	Type      TransactionType
	Date      time.Time // Only contains date; we assume Transactions come in chronological order
}

// Summary represents a stock's OHLC and previous price data
type Summary struct {
	StockCode string    `json:"stock_code"`
	Date      time.Time `json:"date"`
	Prev      int64     `json:"prev"`
	Open      int64     `json:"open"`
	High      int64     `json:"high"`
	Low       int64     `json:"low"`
	Close     int64     `json:"close"`
	Volume    int64     `json:"volume"`
	Value     int64     `json:"value"`
	Average   int64     `json:"average"`
}

// ApplyTransaction returns stockSummary with updated data based on given transaction
// Assumption: TypeA is only used to set Prev price when the Quantity is 0
func (summary Summary) ApplyTransaction(transaction Transaction) (bool, Summary) {
	var (
		updatedSummary = summary
		isUpdated      bool
	)

	if updatedSummary == (Summary{}) {
		updatedSummary.StockCode = transaction.StockCode
		updatedSummary.Date = transaction.Date
	}

	switch transaction.Type {
	case TransactionTypeA:
		if transaction.Quantity == 0 {
			updatedSummary.Prev = transaction.Price
		}
	case TransactionTypeE, TransactionTypeP:
		updatedSummary.Value = summary.Value + (transaction.Quantity * transaction.Price)
		updatedSummary.Volume = summary.Volume + transaction.Quantity
		updatedSummary.Average = (updatedSummary.Value / updatedSummary.Volume)
		fallthrough
	default:
		// Open
		if summary.Open == 0 && transaction.Quantity > 0 {
			updatedSummary.Open = transaction.Price
		}

		// High
		if summary.High == 0 || summary.High < transaction.Price {
			updatedSummary.High = transaction.Price
		}

		// Low
		if summary.Low == 0 || summary.Low > transaction.Price {
			updatedSummary.Low = transaction.Price
		}

		// Close
		updatedSummary.Close = transaction.Price
	}

	isUpdated = summary != updatedSummary
	return isUpdated, updatedSummary
}
