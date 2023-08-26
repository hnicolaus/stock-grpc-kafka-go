package model

import (
	"fmt"
	"strconv"
	"time"
)

type Type string

const (
	TypeA         Type = "A"
	TypeP         Type = "P"
	TypeE         Type = "E"
	TypeUndefined Type = ""
)

type Transaction struct {
	Price     int64
	Quantity  int64
	StockCode string
	Type      Type
	Date      time.Time // Does not contain time as we assume ChangeRecord comes in chronological order regardless of time.
}

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
	quantity := i.Quantity
	// if inputType == TypeE {
	// 	price = i.ExecutionPrice
	// 	quantity = i.ExecutedQuantity
	// }

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

// Summary represents a stock's OHLC and previous price data
type Summary struct {
	StockCode string    `json:"stock_code"`
	Date      time.Time `json:"date"`
	Prev      int64     `json:"prev"`
	Open      int64     `json:"open"`
	High      int64     `json:"high"`
	Low       int64     `json:"low"`
	Close     int64     `json:"close"`
	Volume    int64     `json:"volume"`
	Value     int64     `json:"value"`
	Average   int64     `json:"average"`
}

// ApplyTransaction returns stockSummary with updated data based on given transaction
// Assumption: TypeA is only used to set Prev price when the Quantity is 0
func (summary Summary) ApplyTransaction(transaction Transaction) (bool, Summary) {
	var (
		updatedSummary Summary = summary
		isUpdated      bool    = false
	)

	switch transaction.Type {
	case TypeA:
		if transaction.Quantity == 0 {
			updatedSummary.Prev = transaction.Price
		}
	case TypeE, TypeP:
		updatedSummary.Value = summary.Value + (transaction.Quantity * transaction.Price)
		updatedSummary.Volume = summary.Volume + transaction.Quantity
		updatedSummary.Average = (updatedSummary.Value / updatedSummary.Volume)
		fallthrough
	default:
		// Open
		if summary.Open == 0 && transaction.Quantity > 0 {
			updatedSummary.Open = transaction.Price
		}

		// High
		if summary.High == 0 || summary.High < transaction.Price {
			updatedSummary.High = transaction.Price
		}

		// Low
		if summary.Low == 0 || summary.Low > transaction.Price {
			updatedSummary.Low = transaction.Price
		}

		// Close
		updatedSummary.Close = transaction.Price
	}

	isUpdated = summary != updatedSummary
	return isUpdated, updatedSummary
}
