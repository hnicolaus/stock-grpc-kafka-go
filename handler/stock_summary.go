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

	stockSummaries, err := h.stockUsecase.GetStockSummary(ctx, request)
	if err != nil {
		return &proto.GetStockSummaryResponse{}, err
	}

	response := convertResponseToProto(stockSummaries)
	return response, nil
}

func convertProtoToRequest(req *proto.GetStockSummaryRequest) (model.GetStockSummaryRequest, error) {
	stockCode := req.GetStockCode()
	if stockCode == "" {
		return model.GetStockSummaryRequest{}, errors.New("stockCode cannot be empty")
	}

	toDateString := req.GetToDate()
	if toDateString == "" {
		return model.GetStockSummaryRequest{}, errors.New("toDate cannot be empty")
	}
	toDate, err := time.Parse(stockSummaryDateFmt, toDateString)
	if err != nil {
		return model.GetStockSummaryRequest{}, errors.New("invalid toDate format; please input string with format yyyy-mm-dd")
	}

	fromDateString := req.GetFromDate()
	if fromDateString == "" {
		return model.GetStockSummaryRequest{}, errors.New("fromDateString cannot be empty")
	}
	fromDate, err := time.Parse(stockSummaryDateFmt, fromDateString)
	if err != nil {
		return model.GetStockSummaryRequest{}, errors.New("invalid fromDate format, please input string with format yyyy-mm-dd")
	}

	if fromDate.After(toDate) {
		return model.GetStockSummaryRequest{}, errors.New("toDate must be before or equal to fromDate")
	}

	return model.GetStockSummaryRequest{
		StockCode: stockCode,
		FromDate:  fromDate,
		ToDate:    toDate,
	}, nil
}

func convertResponseToProto(response []model.Summary) *proto.GetStockSummaryResponse {
	result := &proto.GetStockSummaryResponse{
		Result: []*proto.StockSummary{},
	}
	for _, stockSummary := range response {
		result.Result = append(result.Result, &proto.StockSummary{
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
		})
	}

	return result
}
