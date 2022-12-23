package exchange

import (
	"context"

	"github.com/project-code-io/crypto-trading-bot-go/pair"
)

// Binance is an exchange that communicates with the binance exchange,
// configuration with either the .us or the .com domain can be
// done in the constructor.
type Binance struct {
	BaseURL string
}

// NewBinance acts as the default constructor for the Binance exchange type.
// This method takes a BinanceDomain, which is used for specifying either the
// .us domain or the .com domain.
func NewBinance(domain BinanceDomain) *Binance {
	e := &Binance{
		BaseURL: domain.baseURL(),
	}

	return e
}

// GetLastPrice obtains the last price for the pair on binance.
func (e *Binance) GetLastPrice(ctx context.Context, p pair.Pair) (string, error) {
	return "0.00", nil
}

// BinanceDomain is an enum type that is used to specify which domain the
// Binance exchange client should interface with.
type BinanceDomain int

const (
	// BinanceDomainUS specifies the binance.us domain.
	BinanceDomainUS BinanceDomain = 1

	// BinanceDomainCom specifies the binance.com domain.
	BinanceDomainCom = 2
)

func (d BinanceDomain) baseURL() string {
	switch d {
	case BinanceDomainUS:
		return "https://api.binance.us"
	case BinanceDomainCom:
		return "https://api.binance.com"
	default:
		return ""
	}
}
