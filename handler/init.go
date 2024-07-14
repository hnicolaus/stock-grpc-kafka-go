/*
	Hans Nicolaus
	29 Aug 2023
*/

package handler

import (
	"context"

	"stock/model"
	"stock/proto"
)

//go:generate mockgen -source=./init.go -destination=./_mock/stock_summary_mock.go -package=mock
type StockUsecase interface {
	UpdateStockSummary(ctx context.Context, transaction model.Transaction) error
	GetStockSummary(ctx context.Context, request model.GetStockSummaryRequest) ([]model.Summary, error)
}

type Handler struct {
	proto.UnimplementedStockServer
	stockUsecase StockUsecase
}

func New(stockUsecase StockUsecase) *Handler {
	return &Handler{
		stockUsecase: stockUsecase,
	}
}
