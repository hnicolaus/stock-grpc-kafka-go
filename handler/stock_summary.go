package handler

import (
	"context"
	"errors"
	"time"

	"bibit.id/challenge/model"
	"bibit.id/challenge/proto"
)

const (
	stockSummaryDateFmt = "2006-01-02"
)

func (h *Handler) GetStockSummary(ctx context.Context, req *proto.GetStockSummaryRequest) (*proto.GetStockSummaryResponse, error) {
	request, err := convertProtoToRequest(req)
	if err != nil {
		return &proto.GetStockSummaryResponse{}, err
	}

	stockSummary, err := h.stockUsecase.GetStockSummary(ctx, request)
	if err != nil {
		return &proto.GetStockSummaryResponse{}, err
	}

	return &proto.GetStockSummaryResponse{
		StockCode: stockSummary.StockCode,
		Date:      stockSummary.Date.Format(stockSummaryDateFmt),
		Prev:      stockSummary.Prev,
		Open:      stockSummary.Open,
		High:      stockSummary.High,
		Low:       stockSummary.Low,
		Close:     stockSummary.Close,
		Volume:    stockSummary.Volume,
		Value:     stockSummary.Value,
		Average:   stockSummary.Average,
	}, nil
}

func convertProtoToRequest(req *proto.GetStockSummaryRequest) (model.GetStockSummaryRequest, error) {
	stockCode := req.GetStockCode()
	if stockCode == "" {
		return model.GetStockSummaryRequest{}, errors.New("stockCode cannot be empty")
	}

	dateString := req.GetDate()
	if dateString == "" {
		return model.GetStockSummaryRequest{}, errors.New("date cannot be empty")
	}
	date, err := time.Parse(stockSummaryDateFmt, dateString)
	if err != nil {
		return model.GetStockSummaryRequest{}, errors.New("date format is invalid")
	}

	return model.GetStockSummaryRequest{
		StockCode: stockCode,
		Date:      date,
	}, nil
}
