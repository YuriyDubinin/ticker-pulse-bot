package quotes

import "time"

type QuoteInfo struct {
	QuoteID   string
	Ticker    string
	Label     string
	MinPrice  float64
	MaxPrice  float64
	UpdatedAt time.Time
}

// Временное решение до подъема базы
var Quotes = []QuoteInfo{
	{QuoteID: "bitcoin", Ticker: "BTC", Label: "Bitcoin"},
	{QuoteID: "ethereum", Ticker: "ETH", Label: "Ethereum"},
	{QuoteID: "tether", Ticker: "USDT", Label: "Tether"},
	{QuoteID: "binancecoin", Ticker: "BNB", Label: "Binance Coin"},
	{QuoteID: "usd-coin", Ticker: "USDC", Label: "USD Coin"},
	{QuoteID: "ripple", Ticker: "XRP", Label: "Ripple"},
	{QuoteID: "cardano", Ticker: "ADA", Label: "Cardano"},
	{QuoteID: "dogecoin", Ticker: "DOGE", Label: "Dogecoin"},
	{QuoteID: "solana", Ticker: "SOL", Label: "Solana"},
	{QuoteID: "the-open-network", Ticker: "TON", Label: "The Open Network"},
}
