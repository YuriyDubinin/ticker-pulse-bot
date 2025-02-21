package data_formatter

import (
	"errors"
	"fmt"
	"math"
	"time"
)

type QuoteInfo struct {
	CryptoID  string
	Ticker    string
	Label     string
	MinPrice  float64
	MaxPrice  float64
	UpdatedAt time.Time
}

type HistoricalQuoteMinMaxData struct {
	MinPrice float64
	MaxPrice float64
}

// Специально сделан для обработки данных полученных с помощью CalculateQuoteLimitValues,
// пщиск мин/макс показателя для тикера из предоставленных исторических данных
func CalculateHistoricalMinMax(data map[string]interface{}, field string) (HistoricalQuoteMinMaxData, error) {
	var result HistoricalQuoteMinMaxData
	result.MinPrice = math.MaxFloat64
	result.MaxPrice = -math.MaxFloat64

	entries, ok := data[field]
	if !ok {
		return result, fmt.Errorf("field '%s' not found in data", field)
	}

	entriesSlice, ok := entries.([]interface{})
	if !ok {
		return result, errors.New("invalid data format: expected slice of arrays")
	}

	for _, entry := range entriesSlice {
		entryArr, ok := entry.([]interface{})
		if !ok || len(entryArr) != 2 {
			return result, errors.New("invalid data format: expected array of [timestamp, value]")
		}

		value, ok := entryArr[1].(float64)
		if !ok {
			return result, errors.New("invalid value format")
		}

		if value < result.MinPrice {
			result.MinPrice = value
		}
		if value > result.MaxPrice {
			result.MaxPrice = value
		}
	}

	return result, nil
}

// func CombainQuotesMap() {

// }
