/*
	Hans Nicolaus
	29 Aug 2023
*/

package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"bibit.id/challenge/model"
	mock "bibit.id/challenge/usecase/_mock"
	"github.com/golang/mock/gomock"
)

func Test_GetStockSummary(t *testing.T) {
	type args struct {
		ctx   context.Context
		input model.GetStockSummaryRequest
	}
	type fields struct {
		stockRepo func(ctrl *gomock.Controller) StockRepo
	}
	tests := []struct {
		name   string
		args   args
		fields fields

		wantResponse []model.Summary
		wantErr      bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				input: model.GetStockSummaryRequest{
					StockCode: "BBCA",
					FromDate:  time.Time{}.AddDate(0, 0, 1),
					ToDate:    time.Time{}.AddDate(0, 0, 2),
				},
			},
			fields: fields{
				stockRepo: func(ctrl *gomock.Controller) StockRepo {
					m := mock.NewMockStockRepo(ctrl)

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
			wantResponse: []model.Summary{
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
			},
		},
		{
			name: "success-no-result",
			args: args{
				ctx: context.Background(),
				input: model.GetStockSummaryRequest{
					StockCode: "BBCA",
					FromDate:  time.Time{}.AddDate(0, 0, 1),
					ToDate:    time.Time{}.AddDate(0, 0, 2),
				},
			},
			fields: fields{
				stockRepo: func(ctrl *gomock.Controller) StockRepo {
					m := mock.NewMockStockRepo(ctrl)

					m.EXPECT().GetStockSummary(gomock.Any(), model.GetStockSummaryRequest{
						StockCode: "BBCA",
						FromDate:  time.Time{}.AddDate(0, 0, 1),
						ToDate:    time.Time{}.AddDate(0, 0, 2),
					}).Return([]model.Summary{}, nil)

					return m
				},
			},
			wantResponse: []model.Summary{},
		},
		{
			name: "error-get-stock-summary",
			args: args{
				ctx: context.Background(),
				input: model.GetStockSummaryRequest{
					StockCode: "BBCA",
					FromDate:  time.Time{}.AddDate(0, 0, 1),
					ToDate:    time.Time{}.AddDate(0, 0, 2),
				},
			},
			fields: fields{
				stockRepo: func(ctrl *gomock.Controller) StockRepo {
					m := mock.NewMockStockRepo(ctrl)

					m.EXPECT().GetStockSummary(gomock.Any(), model.GetStockSummaryRequest{
						StockCode: "BBCA",
						FromDate:  time.Time{}.AddDate(0, 0, 1),
						ToDate:    time.Time{}.AddDate(0, 0, 2),
					}).Return([]model.Summary{}, errors.New("error-get-stock-summary"))

					return m
				},
			},
			wantResponse: []model.Summary{},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			usecase := &Usecase{
				stockRepo: tt.fields.stockRepo(ctrl),
			}

			gotResponse, err := usecase.GetStockSummary(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.GetStockSummary() err = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("usecase.GetStockSummary() gotResponse = %v, wantResponse %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func Test_UpdateStockSummary(t *testing.T) {
	type args struct {
		ctx   context.Context
		input model.Transaction
	}
	type fields struct {
		stockRepo func(ctrl *gomock.Controller) StockRepo
	}
	tests := []struct {
		name   string
		args   args
		fields fields

		wantErr bool
	}{
		{
			name: "success-type-a-qty-0-no-prev-updates-prev",
			args: args{
				ctx: context.Background(),
				input: model.Transaction{
					StockCode: "BBCA",
					Price:     8000,
					Quantity:  0,
					Type:      model.TypeA,
					Date:      time.Time{}.AddDate(0, 0, 1),
				},
			},
			fields: fields{
				stockRepo: func(ctrl *gomock.Controller) StockRepo {
					m := mock.NewMockStockRepo(ctrl)

					m.EXPECT().GetStockSummary(gomock.Any(), model.GetStockSummaryRequest{
						StockCode: "BBCA",
						FromDate:  time.Time{}.AddDate(0, 0, 1),
						ToDate:    time.Time{}.AddDate(0, 0, 1),
					}).Return([]model.Summary{}, nil)

					m.EXPECT().UpdateStockSummary(gomock.Any(), model.Summary{
						StockCode: "BBCA",
						Date:      time.Time{}.AddDate(0, 0, 1),
						Prev:      8000,
					}).Return(nil)

					return m
				},
			},
		},
		{
			name: "success-type-a-qty-0-existing-prev-updates-prev",
			args: args{
				ctx: context.Background(),
				input: model.Transaction{
					StockCode: "BBCA",
					Price:     10000,
					Quantity:  0,
					Type:      model.TypeA,
					Date:      time.Time{}.AddDate(0, 0, 1),
				},
			},
			fields: fields{
				stockRepo: func(ctrl *gomock.Controller) StockRepo {
					m := mock.NewMockStockRepo(ctrl)

					m.EXPECT().GetStockSummary(gomock.Any(), model.GetStockSummaryRequest{
						StockCode: "BBCA",
						FromDate:  time.Time{}.AddDate(0, 0, 1),
						ToDate:    time.Time{}.AddDate(0, 0, 1),
					}).Return([]model.Summary{
						{
							StockCode: "BBCA",
							Date:      time.Time{}.AddDate(0, 0, 1),
							Prev:      8000,
						},
					}, nil)

					m.EXPECT().UpdateStockSummary(gomock.Any(), model.Summary{
						StockCode: "BBCA",
						Date:      time.Time{}.AddDate(0, 0, 1),
						Prev:      10000,
					}).Return(nil)

					return m
				},
			},
		},
		{
			name: "success-type-p-updates-data-1",
			args: args{
				ctx: context.Background(),
				input: model.Transaction{
					StockCode: "BBCA",
					Price:     8050,
					Quantity:  100,
					Type:      model.TypeP,
					Date:      time.Time{}.AddDate(0, 0, 1),
				},
			},
			fields: fields{
				stockRepo: func(ctrl *gomock.Controller) StockRepo {
					m := mock.NewMockStockRepo(ctrl)

					m.EXPECT().GetStockSummary(gomock.Any(), model.GetStockSummaryRequest{
						StockCode: "BBCA",
						FromDate:  time.Time{}.AddDate(0, 0, 1),
						ToDate:    time.Time{}.AddDate(0, 0, 1),
					}).Return([]model.Summary{
						{
							StockCode: "BBCA",
							Date:      time.Time{}.AddDate(0, 0, 1),
							Prev:      8000,
						},
					}, nil)

					m.EXPECT().UpdateStockSummary(gomock.Any(), model.Summary{
						StockCode: "BBCA",
						Date:      time.Time{}.AddDate(0, 0, 1),
						Prev:      8000,
						Open:      8050,
						High:      8050,
						Low:       8050,
						Close:     8050,
						Volume:    100,
						Value:     805000,
						Average:   8050,
					}).Return(nil)

					return m
				},
			},
		},
		{
			name: "success-type-p-updates-data-2",
			args: args{
				ctx: context.Background(),
				input: model.Transaction{
					StockCode: "BBCA",
					Price:     7950,
					Quantity:  500,
					Type:      model.TypeP,
					Date:      time.Time{}.AddDate(0, 0, 1),
				},
			},
			fields: fields{
				stockRepo: func(ctrl *gomock.Controller) StockRepo {
					m := mock.NewMockStockRepo(ctrl)

					m.EXPECT().GetStockSummary(gomock.Any(), model.GetStockSummaryRequest{
						StockCode: "BBCA",
						FromDate:  time.Time{}.AddDate(0, 0, 1),
						ToDate:    time.Time{}.AddDate(0, 0, 1),
					}).Return([]model.Summary{
						{
							StockCode: "BBCA",
							Date:      time.Time{}.AddDate(0, 0, 1),
							Prev:      8000,
							Open:      8050,
							High:      8050,
							Low:       8050,
							Close:     8050,
							Volume:    100,
							Value:     805000,
							Average:   8050,
						},
					}, nil)

					m.EXPECT().UpdateStockSummary(gomock.Any(), model.Summary{
						StockCode: "BBCA",
						Date:      time.Time{}.AddDate(0, 0, 1),
						Prev:      8000,
						Open:      8050,
						High:      8050,
						Low:       7950,
						Close:     7950,
						Volume:    600,
						Value:     4780000,
						Average:   7966,
					}).Return(nil)

					return m
				},
			},
		},
		{
			name: "success-type-a-non0-qty-no-update-1",
			args: args{
				ctx: context.Background(),
				input: model.Transaction{
					StockCode: "BBCA",
					Price:     8150,
					Quantity:  200,
					Type:      model.TypeA,
					Date:      time.Time{}.AddDate(0, 0, 1),
				},
			},
			fields: fields{
				stockRepo: func(ctrl *gomock.Controller) StockRepo {
					m := mock.NewMockStockRepo(ctrl)

					m.EXPECT().GetStockSummary(gomock.Any(), model.GetStockSummaryRequest{
						StockCode: "BBCA",
						FromDate:  time.Time{}.AddDate(0, 0, 1),
						ToDate:    time.Time{}.AddDate(0, 0, 1),
					}).Return([]model.Summary{
						{
							StockCode: "BBCA",
							Date:      time.Time{}.AddDate(0, 0, 1),
							Prev:      8000,
							Open:      8050,
							High:      8050,
							Low:       7950,
							Close:     7950,
							Volume:    600,
							Value:     4780000,
							Average:   7966,
						},
					}, nil)

					return m
				},
			},
		},
		{
			name: "success-type-e-updates-data",
			args: args{
				ctx: context.Background(),
				input: model.Transaction{
					StockCode: "BBCA",
					Price:     8100,
					Quantity:  300,
					Type:      model.TypeE,
					Date:      time.Time{}.AddDate(0, 0, 1),
				},
			},
			fields: fields{
				stockRepo: func(ctrl *gomock.Controller) StockRepo {
					m := mock.NewMockStockRepo(ctrl)

					m.EXPECT().GetStockSummary(gomock.Any(), model.GetStockSummaryRequest{
						StockCode: "BBCA",
						FromDate:  time.Time{}.AddDate(0, 0, 1),
						ToDate:    time.Time{}.AddDate(0, 0, 1),
					}).Return([]model.Summary{
						{
							StockCode: "BBCA",
							Date:      time.Time{}.AddDate(0, 0, 1),
							Prev:      8000,
							Open:      8050,
							High:      8050,
							Low:       7950,
							Close:     7950,
							Volume:    600,
							Value:     4780000,
							Average:   7966,
						},
					}, nil)

					m.EXPECT().UpdateStockSummary(gomock.Any(), model.Summary{
						StockCode: "BBCA",
						Date:      time.Time{}.AddDate(0, 0, 1),
						Prev:      8000,
						Open:      8050,
						High:      8100,
						Low:       7950,
						Close:     8100,
						Volume:    900,
						Value:     7210000,
						Average:   8011,
					}).Return(nil)

					return m
				},
			},
		},
		{
			name: "success-type-a-non0-qty-no-update-2",
			args: args{
				ctx: context.Background(),
				input: model.Transaction{
					StockCode: "BBCA",
					Price:     8200,
					Quantity:  100,
					Type:      model.TypeA,
					Date:      time.Time{}.AddDate(0, 0, 1),
				},
			},
			fields: fields{
				stockRepo: func(ctrl *gomock.Controller) StockRepo {
					m := mock.NewMockStockRepo(ctrl)

					m.EXPECT().GetStockSummary(gomock.Any(), model.GetStockSummaryRequest{
						StockCode: "BBCA",
						FromDate:  time.Time{}.AddDate(0, 0, 1),
						ToDate:    time.Time{}.AddDate(0, 0, 1),
					}).Return([]model.Summary{
						{
							StockCode: "BBCA",
							Date:      time.Time{}.AddDate(0, 0, 1),
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			usecase := &Usecase{
				stockRepo: tt.fields.stockRepo(ctrl),
			}

			err := usecase.UpdateStockSummary(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.UpdateStockSummary() err = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
