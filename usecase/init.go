package usecase

import (
	"time"

	"bibit.id/challenge/model"
)

type StockRepo interface {
	GetStockSummary(stockCode string, date time.Time) (summary model.Summary, err error)
	SetStockSummary(stockCode string, date time.Time, summary model.Summary) (err error)
}

type Usecase struct {
	stockRepo StockRepo
}

func New(stockRepo StockRepo) *Usecase {
	return &Usecase{
		stockRepo: stockRepo,
	}
}
