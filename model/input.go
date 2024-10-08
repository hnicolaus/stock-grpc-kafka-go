/*
	Hans Nicolaus
	29 Aug 2023
*/

package model

import (
	"fmt"
	"strconv"
	"time"
)

type KafkaTransaction struct {
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

func (i *KafkaTransaction) ToTransaction() (Transaction, error) {
	inputType := convertToType(i.Type)
	if inputType == TransactionTypeUndefined {
		return Transaction{}, fmt.Errorf("invalid transaction type %s", i.Type)
	}

	price := i.Price
	if price == "" {
		price = i.ExecutionPrice
	}

	quantity := i.Quantity
	if quantity == "" {
		quantity = i.ExecutedQuantity
	}

	inputPrice, err := strconv.ParseInt(price, 10, 64)
	if err != nil {
		return Transaction{}, err
	}

	inputQuantity := int64(0)
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

	timestamp, err := time.Parse(timeFormat, orderNumber) // assume order number contains date in the "yyyyMMdd" format
	if err != nil {
		return time.Time{}, err
	}

	return timestamp, nil
}

func convertToType(t string) TransactionType {
	switch t {
	case "A":
		return TransactionTypeA
	case "P":
		return TransactionTypeP
	case "E":
		return TransactionTypeE
	default:
		return TransactionTypeUndefined
	}
}

type GetStockSummaryRequest struct {
	StockCode string
	FromDate  time.Time
	ToDate    time.Time
}
