package model

import (
	"fmt"
	"strconv"
	"time"
)

type Input struct {
	Type             string `json:"type"`
	OrderBook        string `json:"order_book,omitempty"`
	OrderNumber      string `json:"order_number,omitempty"`
	OrderVerb        string `json:"order_verb,omitempty"`
	Quantity         string `json:"quantity,omitempty"`
	Price            string `json:"price,omitempty"`
	StockCode        string `json:"stock_code,omitempty"`
	ExecutedQuantity string `json:"executed_quantity,omitempty"`
	ExecutionPrice   string `json:"execution_price,omitempty"`
}

func (i *Input) ToTransaction() (Transaction, error) {
	var (
		inputType                 Type
		inputPrice, inputQuantity int64

		err error
	)

	inputType = convertToType(i.Type)
	if inputType == TypeUndefined {
		return Transaction{}, fmt.Errorf("invalid type: %s", i.Type)
	}

	price := i.Price
	if price == "" {
		price = i.ExecutionPrice
	}

	quantity := i.Quantity
	if quantity == "" {
		quantity = i.ExecutedQuantity
	}

	inputPrice, err = strconv.ParseInt(price, 10, 64)
	if err != nil {
		return Transaction{}, err
	}

	if quantity != "" {
		inputQuantity, err = strconv.ParseInt(quantity, 10, 64)
		if err != nil {
			return Transaction{}, err
		}
	}

	// Assume that OrderNumber contains timestamp in the format of yyyyMMddHHmmss
	timestamp, err := getDateFromOrderNumber(i.OrderNumber)
	if err != nil {
		return Transaction{}, err
	}

	return Transaction{
		Type:      inputType,
		Price:     inputPrice,
		Quantity:  inputQuantity,
		StockCode: i.StockCode,
		Date:      timestamp,
	}, nil
}

func getDateFromOrderNumber(orderNumber string) (time.Time, error) {
	timeFormat := "20060102"
	if len(orderNumber) > len(timeFormat) {
		orderNumber = orderNumber[:len(timeFormat)]
	}

	timestamp, err := time.Parse(timeFormat, orderNumber) // Assuming the format is "yyyyMMdd"
	if err != nil {
		return time.Time{}, err
	}

	return timestamp, nil
}

func convertToType(t string) Type {
	switch t {
	case "A":
		return TypeA
	case "P":
		return TypeP
	case "E":
		return TypeE
	default:
		return TypeUndefined
	}
}

type GetStockSummaryRequest struct {
	StockCode string
	Date      time.Time
}
