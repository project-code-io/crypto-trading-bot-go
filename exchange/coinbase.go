package exchange

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/project-code-io/crypto-trading-bot-go/order"
	"github.com/project-code-io/crypto-trading-bot-go/trading"
)

// Coinbase represents a client that is able to talk to the coinbase advanced
// trading api.
type Coinbase struct {
	APIKey    string
	APISecret string
}

const coinbaseBaseURL = "https://api.coinbase.com"

// NewCoinbase acts as the default constructor for the Coinbase exchange type.
// This method will attempt to load authentication credentials from the
// environment, returning an error if any are missing.
func NewCoinbase() (*Coinbase, error) {
	key, exists := os.LookupEnv("COINBASE_API_KEY")
	if !exists {
		return nil, ErrAPIKeyNotSet
	}

	secret, exists := os.LookupEnv("COINBASE_API_SECRET")
	if !exists {
		return nil, ErrAPISecretNotSet
	}

	e := &Coinbase{
		APIKey:    key,
		APISecret: secret,
	}

	return e, nil
}

func (e *Coinbase) getBody(r *http.Request) (string, error) {
	if r.Body == nil {
		return "", nil
	}

	bodyReader, err := r.GetBody()
	if err != nil {
		return "", fmt.Errorf("get body of req: %w", err)
	}

	bodyData, err := io.ReadAll(bodyReader)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	return string(bodyData), nil
}

func (e *Coinbase) doRequest(r *http.Request) (*http.Response, error) {
	timestamp := time.Now().Unix()

	body, err := e.getBody(r)
	if err != nil {
		return nil, fmt.Errorf("get body: %w", err)
	}

	payload := fmt.Sprintf("%d%s%s%s", timestamp, r.Method, r.URL.Path, body)

	hash := hmac.New(sha256.New, []byte(e.APISecret))
	hash.Write([]byte(payload))
	sig := hex.EncodeToString(hash.Sum(nil))

	r.Header.Add("CB-ACCESS-KEY", e.APIKey)
	r.Header.Add("CB-ACCESS-SIGN", sig)
	r.Header.Add("CB-ACCESS-TIMESTAMP", strconv.Itoa(int(timestamp)))

	return http.DefaultClient.Do(r)
}

func (e *Coinbase) convertPairValue(p trading.Pair) (string, error) {
	switch p {
	case trading.BTCUSD:
		return "BTC-USD", nil
	case trading.ETHUSD:
		return "ETH-USD", nil
	default:
		return "", ErrMissingPair
	}
}

// GetLastPrice obtains the last price for the pair on binance.
func (e *Coinbase) GetLastPrice(ctx context.Context, p trading.Pair) (string, error) {
	type priceResponse struct {
		Price string `json:"price"`
	}

	pairVal, err := e.convertPairValue(p)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/api/v3/brokerage/products/%s", coinbaseBaseURL, pairVal)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	res, err := e.doRequest(req)
	if err != nil {
		return "", err
	}

	var response priceResponse

	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("decode json: %w", err)
	}

	return response.Price, nil
}

func (e *Coinbase) CreateLimitOrder(ctx context.Context, order order.Limit) (Order, error) {
	return Order{}, nil
}

func (e *Coinbase) CancelOrders(ctx context.Context, orderID ...string) error {
	return nil
}

func (e *Coinbase) ListOpenOrders(ctx context.Context) ([]Order, error) {
	return nil, nil
}

func (e *Coinbase) GetBalance(ctx context.Context, asset trading.Asset) (int64, error) {
	return 0, nil
}
