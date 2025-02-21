package crypto_fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type CryptoFetcher struct {
	baseURL string
}

func NewCryptoFetcher() *CryptoFetcher {
	return &CryptoFetcher{
		baseURL: "https://api.coingecko.com/api/v3",
	}
}

func (cf *CryptoFetcher) FetchCoingeckoQuotesRate(cryptoID string, currency string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=%s", cf.baseURL, cryptoID, currency)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[TICKER-PULSE-BOT]: [CRYPTO-FETCHER]: Ошибка запроса к Coingecko: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("неудачный запрос: %s", resp.Status)
		log.Printf("[TICKER-PULSE-BOT]: [CRYPTO-FETCHER]: %v", err)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[TICKER-PULSE-BOT]: [CRYPTO-FETCHER]: Ошибка чтения тела ответа: %v", err)
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("[TICKER-PULSE-BOT]: [CRYPTO-FETCHER]: Ошибка парсинга JSON: %v", err)
		return nil, err
	}

	return result, nil
}

func (cf *CryptoFetcher) FetchCoingeckoHistoricalData(cryptoID string, period int) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/coins/%s/market_chart?vs_currency=usd&days=%d&interval=daily", cf.baseURL, cryptoID, period)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[TICKER-PULSE-BOT]: [CRYPTO-FETCHER]: Ошибка запроса к Coingecko: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("неудачный запрос: %s", resp.Status)
		log.Printf("[TICKER-PULSE-BOT]: [CRYPTO-FETCHER]: %v", err)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[TICKER-PULSE-BOT]: [CRYPTO-FETCHER]: Ошибка чтения тела ответа: %v", err)
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("[TICKER-PULSE-BOT]: [CRYPTO-FETCHER]: Ошибка парсинга JSON: %v", err)
		return nil, err
	}

	return result, nil
}
