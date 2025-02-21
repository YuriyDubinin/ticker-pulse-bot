package telegram_message_formatter

import (
	"fmt"
	"sort"
	"strings"
)

func ConvertCryptoDataToString(data map[string]interface{}) string {
	var (
		result strings.Builder
		keys   = make([]string, 0, len(data))
		prices = make(map[string]float64, len(data))
	)

	// Преобразование данных и сбор ключей
	for key, value := range data {
		if innerMap, ok := value.(map[string]interface{}); ok {
			if price, ok := innerMap["usd"].(float64); ok {
				keys = append(keys, key)
				prices[key] = price
			}
		}
	}

	sort.Strings(keys)

	for _, key := range keys {
		result.WriteString(fmt.Sprintf("%-20s %15.2f  $\n", key, prices[key]))
	}

	return result.String()
}
