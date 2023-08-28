package usecase

import (
	"context"

	"bibit.id/challenge/model"
)

//go:generate mockgen -source=./init.go -destination=./_mock/stock_summary_mock.go -package=mock
type StockRepo interface {
	GetStockSummary(ctx context.Context, request model.GetStockSummaryRequest) (result []model.Summary, err error)
	UpdateStockSummary(ctx context.Context, stockSummary model.Summary) (err error)
}

type Usecase struct {
	stockRepo StockRepo
}

func New(stockRepo StockRepo) *Usecase {
	return &Usecase{
		stockRepo: stockRepo,
	}
}
