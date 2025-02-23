package data_formatter

import (
	"errors"
	"fmt"
	"math"
	"strings"
	quotes "ticker-pulse-bot/internal/pkg/quotes"
)

type HistoricalQuoteMinMaxData struct {
	MinPrice float64
	MaxPrice float64
}

// Специально сделан для обработки данных полученных с помощью CalculateQuoteLimitValues,
// поиск мин/макс показателя для тикера из предоставленных исторических данных
func CalculateHistoricalMinMax(data map[string]interface{}, field string) (HistoricalQuoteMinMaxData, error) {
	var result HistoricalQuoteMinMaxData
	result.MinPrice = math.MaxFloat64
	result.MaxPrice = -math.MaxFloat64

	entries, ok := data[field]
	if !ok {
		return result, fmt.Errorf("[CalculateHistoricalMinMax]: field '%s' not found in data", field)
	}

	entriesSlice, ok := entries.([]interface{})
	if !ok {
		return result, errors.New("[CalculateHistoricalMinMax]: invalid data format: expected slice of arrays")
	}

	for _, entry := range entriesSlice {
		entryArr, ok := entry.([]interface{})
		if !ok || len(entryArr) != 2 {
			return result, errors.New("[CalculateHistoricalMinMax]: invalid data format: expected array of [timestamp, value]")
		}

		value, ok := entryArr[1].(float64)
		if !ok {
			return result, errors.New("[CalculateHistoricalMinMax]: invalid value format")
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

func FormatQuotesToString() (string, error) {
	if len(quotes.Quotes) == 0 {
		return "", fmt.Errorf("[FormatQuotesToString]: quotes list is empty")
	}

	var tickers []string
	for _, q := range quotes.Quotes {
		if q.QuoteID == "" {
			return "", fmt.Errorf("[FormatQuotesToString]: found empty QuoteID in quotes")
		}
		tickers = append(tickers, q.QuoteID)
	}

	return strings.Join(tickers, ","), nil
}
