/*
	Hans Nicolaus
	29 Aug 2023
*/

package handler

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	mock "stock/handler/_mock"
	"stock/model"
	"stock/proto"

	"github.com/golang/mock/gomock"
)

func Test_Handler_GetStockSummary(t *testing.T) {
	type args struct {
		ctx   context.Context
		input *proto.GetStockSummaryRequest
	}
	type fields struct {
		stockUsecase func(ctrl *gomock.Controller) StockUsecase
	}
	tests := []struct {
		name   string
		args   args
		fields fields

		wantResponse *proto.GetStockSummaryResponse
		wantErr      bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				input: &proto.GetStockSummaryRequest{
					StockCode: "BBCA",
					FromDate:  "0001-01-02",
					ToDate:    "0001-01-03",
				},
			},
			fields: fields{
				stockUsecase: func(ctrl *gomock.Controller) StockUsecase {
					m := mock.NewMockStockUsecase(ctrl)

					m.EXPECT().GetStockSummary(gomock.Any(), model.GetStockSummaryRequest{
						StockCode: "BBCA",
						FromDate:  time.Time{}.AddDate(0, 0, 1),
						ToDate:    time.Time{}.AddDate(0, 0, 2),
					}).Return([]model.Summary{
						{
							StockCode: "BBCA",
							Date:      time.Time{}.AddDate(0, 0, 2),
							Prev:      8000,
							Open:      8050,
							High:      8100,
							Low:       7950,
							Close:     8100,
							Volume:    900,
							Value:     7210000,
							Average:   8011,
						},
					}, nil)

					return m
				},
			},
			wantResponse: &proto.GetStockSummaryResponse{
				Result: []*proto.StockSummary{
					{
						StockCode: "BBCA",
						Date:      "0001-01-03",
						Prev:      8000,
						Open:      8050,
						High:      8100,
						Low:       7950,
						Close:     8100,
						Volume:    900,
						Value:     7210000,
						Average:   8011,
					},
				},
			},
		},
		{
			name: "success-no-summary",
			args: args{
				ctx: context.Background(),
				input: &proto.GetStockSummaryRequest{
					StockCode: "BBCA",
					FromDate:  "0001-01-02",
					ToDate:    "0001-01-03",
				},
			},
			fields: fields{
				stockUsecase: func(ctrl *gomock.Controller) StockUsecase {
					m := mock.NewMockStockUsecase(ctrl)

					m.EXPECT().GetStockSummary(gomock.Any(), model.GetStockSummaryRequest{
						StockCode: "BBCA",
						FromDate:  time.Time{}.AddDate(0, 0, 1),
						ToDate:    time.Time{}.AddDate(0, 0, 2),
					}).Return([]model.Summary{}, nil)

					return m
				},
			},
			wantResponse: &proto.GetStockSummaryResponse{
				Result: []*proto.StockSummary{},
			},
		},
		{
			name: "error-get-stock-summary",
			args: args{
				ctx: context.Background(),
				input: &proto.GetStockSummaryRequest{
					StockCode: "BBCA",
					FromDate:  "0001-01-02",
					ToDate:    "0001-01-03",
				},
			},
			fields: fields{
				stockUsecase: func(ctrl *gomock.Controller) StockUsecase {
					m := mock.NewMockStockUsecase(ctrl)

					m.EXPECT().GetStockSummary(gomock.Any(), model.GetStockSummaryRequest{
						StockCode: "BBCA",
						FromDate:  time.Time{}.AddDate(0, 0, 1),
						ToDate:    time.Time{}.AddDate(0, 0, 2),
					}).Return([]model.Summary{
						{
							StockCode: "BBCA",
							Date:      time.Time{}.AddDate(0, 0, 2),
							Prev:      8000,
							Open:      8050,
							High:      8100,
							Low:       7950,
							Close:     8100,
							Volume:    900,
							Value:     7210000,
							Average:   8011,
						},
					}, errors.New("error-get-stock-summary"))

					return m
				},
			},
			wantResponse: &proto.GetStockSummaryResponse{},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			handler := &Handler{
				stockUsecase: tt.fields.stockUsecase(ctrl),
			}

			gotResponse, err := handler.GetStockSummary(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.GetStockSummary() err = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("handler.GetStockSummary() gotResponse = %v, wantResponse %v", gotResponse, tt.wantResponse)
			}
		})
	}
}
