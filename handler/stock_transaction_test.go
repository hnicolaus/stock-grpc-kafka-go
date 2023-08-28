/*
	Hans Nicolaus
	29 Aug 2023
*/

package handler

import (
	"errors"
	"testing"
	"time"

	mock "bibit.id/challenge/handler/_mock"
	"bibit.id/challenge/model"
	"github.com/golang/mock/gomock"
)

func Test_ProcessStockTransaction(t *testing.T) {
	type args struct {
		data []byte
	}
	type fields struct {
		stockUsecase func(ctrl *gomock.Controller) StockUsecase
	}
	tests := []struct {
		name   string
		args   args
		fields fields

		wantErr bool
	}{
		{
			name: "success",
			args: args{
				data: []byte(`{
					"type": "A",
					"quantity": "100",
					"price": "8200",
					"stock_code": "BBCA",
					"order_number": "000101020000073390"
				}`),
			},
			fields: fields{
				stockUsecase: func(ctrl *gomock.Controller) StockUsecase {
					m := mock.NewMockStockUsecase(ctrl)

					m.EXPECT().UpdateStockSummary(gomock.Any(), model.Transaction{
						StockCode: "BBCA",
						Price:     8200,
						Quantity:  100,
						Type:      model.TransactionTypeA,
						Date:      time.Time{}.AddDate(0, 0, 1),
					}).Return(nil)

					return m
				},
			},
		},
		{
			name: "success-type-e-executed-qty-execution-price",
			args: args{
				data: []byte(`{
					"type": "A",
					"executed_quantity": "100",
					"execution_price": "8200",
					"stock_code": "BBCA",
					"order_number": "000101020000073390"
				}`),
			},
			fields: fields{
				stockUsecase: func(ctrl *gomock.Controller) StockUsecase {
					m := mock.NewMockStockUsecase(ctrl)

					m.EXPECT().UpdateStockSummary(gomock.Any(), model.Transaction{
						StockCode: "BBCA",
						Price:     8200,
						Quantity:  100,
						Type:      model.TransactionTypeA,
						Date:      time.Time{}.AddDate(0, 0, 1),
					}).Return(nil)

					return m
				},
			},
		},
		{
			name: "error-update-stock-summary",
			args: args{
				data: []byte(`{
					"type": "A",
					"quantity": "100",
					"price": "8200",
					"stock_code": "BBCA",
					"order_number": "000101020000073390"
				}`),
			},
			fields: fields{
				stockUsecase: func(ctrl *gomock.Controller) StockUsecase {
					m := mock.NewMockStockUsecase(ctrl)

					m.EXPECT().UpdateStockSummary(gomock.Any(), model.Transaction{
						StockCode: "BBCA",
						Price:     8200,
						Quantity:  100,
						Type:      model.TransactionTypeA,
						Date:      time.Time{}.AddDate(0, 0, 1),
					}).Return(errors.New("error-update-stock-summary"))

					return m
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			handler := &Handler{
				stockUsecase: tt.fields.stockUsecase(ctrl),
			}

			err := handler.ProcessStockTransaction(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.ProcessStockTransaction() err = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
