package exchange

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/project-code-io/crypto-trading-bot-go/order"
	"github.com/project-code-io/crypto-trading-bot-go/trading"
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

// ErrBadBinanceDomain describes an error in which the binance domain is not
// correctly set.
var ErrBadBinanceDomain = errors.New(
	"the binance domain does not match the location of the bot",
)

func (e *Binance) convertPairValue(p trading.Pair) (string, error) {
	switch p {
	case trading.BTCUSD:
		return "BTCUSD", nil
	case trading.ETHUSD:
		return "ETHUSD", nil
	default:
		return "", ErrMissingPair
	}
}

// GetLastPrice obtains the last price for the pair on binance.
func (e *Binance) GetLastPrice(ctx context.Context, p trading.Pair) (string, error) {
	type priceResponse struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}

	// https://api.binance.us/api/v3/ticker/price?symbol=ETHUSD
	symbol, err := e.convertPairValue(p)
	if err != nil {
		return "", fmt.Errorf("convert pair value: %w", err)
	}

	url := fmt.Sprintf("%s/api/v3/ticker/price?symbol=%s", e.BaseURL, symbol)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("create new request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("perform request: %w", err)
	}

	const badLocationCode = 451

	if res.StatusCode == badLocationCode {
		return "", ErrBadBinanceDomain
	}

	var data priceResponse

	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("decode response body: %w", err)
	}

	return data.Price, nil
}

// BinanceDomain is an enum type that is used to specify which domain the
// Binance exchange client should interface with.
type BinanceDomain int

const (
	// BinanceDomainUS specifies the binance.us domain.
	BinanceDomainUS BinanceDomain = 1

	// BinanceDomainDotCom specifies the binance.com domain.
	BinanceDomainDotCom = 2
)

func (d BinanceDomain) baseURL() string {
	switch d {
	case BinanceDomainUS:
		return "https://api.binance.us"
	case BinanceDomainDotCom:
		return "https://api.binance.com"
	default:
		return ""
	}
}

func (e *Binance) LimitOrder(ctx context.Context, order order.Limit) (Order, error) {
	return Order{}, nil
}

func (e *Binance) CancelOrder(ctx context.Context, orderID string) error {
	return nil
}

func (e *Binance) ListOpenOrders(ctx context.Context) ([]Order, error) {
	return nil, nil
}

func (e *Binance) GetBalance(ctx context.Context) (int64, error) {
	return 0, nil
}
