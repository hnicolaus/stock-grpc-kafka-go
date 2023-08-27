package handler

import (
	"context"

	"bibit.id/challenge/model"
	"bibit.id/challenge/proto"
)

type StockUsecase interface {
	UpdateStockSummary(ctx context.Context, transaction model.Transaction) error
	GetStockSummary(ctx context.Context, request model.GetStockSummaryRequest) ([]model.Summary, error)
}

type Handler struct {
	proto.UnimplementedBibitServer
	stockUsecase StockUsecase
}

func New(stockUsecase StockUsecase) *Handler {
	return &Handler{
		stockUsecase: stockUsecase,
	}
}
