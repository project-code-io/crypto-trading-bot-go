package exchange

import (
	"bytes"
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
	"strings"
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

func (e *Coinbase) convertAssetvalue(a trading.Asset) (string, error) {
	switch a {
	case trading.BTC:
		return "BTC", nil
	case trading.ETH:
		return "ETH", nil
	case trading.USD:
		return "USD", nil
	}

	return "", ErrMissingAsset
}

func (e *Coinbase) convertSide(s order.Side) string {
	return strings.ToUpper(string(s))
}

func (e *Coinbase) convertToSide(s string) (order.Side, error) {
	switch s {
	case "BUY":
		return order.SideBuy, nil
	case "SELL":
		return order.SideSell, nil
	}

	return order.SideBuy, fmt.Errorf("missing order side conversion")
}

func (e *Coinbase) convertToPair(s string) (trading.Pair, error) {
	switch s {
	case "BTC-USD":
		return trading.BTCUSD, nil
	case "ETH-USD":
		return trading.ETHUSD, nil
	}

	return trading.BTCUSD, fmt.Errorf("missing pair conversion")
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
	type LimitGTC struct {
		BaseSize   string `json: "base_size"`
		LimitPrice string `json: "limit_price"`
		PostOnly   bool   `json: "post_only"`
	}

	type orderConfiguration struct {
		LimitGTC *LimitGTC `json: "limit_limit_gtc"`
	}

	type createOrderRequest struct {
		ClientOrderID      string             `json: "client_order_id"`
		ProductID          string             `json: "product_id"`
		Side               string             `json: "side"`
		OrderConfiguration orderConfiguration `json: "order_configuration"`
	}

	pair, err := e.convertPairValue(order.Pair)
	if err != nil {
		return Order{}, err
	}

	data := createOrderRequest{
		ClientOrderID: order.ClientID,
		ProductID:     pair,
		Side:          e.convertSide(order.Side),
		OrderConfiguration: orderConfiguration{
			LimitGTC: &LimitGTC{
				BaseSize:   order.BaseSize,
				LimitPrice: order.Price,
				PostOnly:   order.PostOnly,
			},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return Order{}, fmt.Errorf("marshal input: %w", err)
	}

	url := fmt.Sprintf("%s/api/v3/brokerage/orders", coinbaseBaseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return Order{}, fmt.Errorf("create request: %w", err)
	}

	res, err := e.doRequest(req)
	if err != nil {
		return Order{}, fmt.Errorf("create request: %w", err)
	}

	type orderResponse struct {
		Success bool   `json: "success"`
		OrderID string `json: "order_id"`
	}

	var orderRes orderResponse

	if err = json.NewDecoder(res.Body).Decode(&orderRes); err != nil {
		return Order{}, fmt.Errorf("decode order response: %w", err)
	}

	if !orderRes.Success {
		return Order{}, ErrOrderFailed
	}

	return Order{
		ID:       orderRes.OrderID,
		Pair:     order.Pair,
		Side:     order.Side,
		ClientID: order.ClientID,
	}, nil
}

func (e *Coinbase) CancelOrders(ctx context.Context, orderID ...string) error {
	// type cancelReqBody struct {
	// 	OrderIDs []string `json: "order_ids"`
	// }

	// bodyData := cancelReqBody{
	// 	OrderIDs: orderID,
	// }

	// data, err := json.Marshal(bodyData)
	// if err != nil {
	// 	return fmt.Errorf("marshal body: %w", err)
	// }

	// url := fmt.Sprintf("%s/api/v3/brokerage/orders/batch_cancel", coinbaseBaseURL)

	// req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	// if err != nil {
	// 	return fmt.Errorf("new request: %w", err)
	// }

	// res, err := e.doRequest(req)
	// if err != nil {
	// 	return fmt.Errorf("perform request: %w", err)
	// }

	// type cancelResponseBody struct {
	// 	Results []cancelResult `json: "result"`
	// }

	// type cancelResult struct {
	// 	Success bool `json: "success"`
	// }

	// var resBody cancelResponseBody

	// if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
	// 	return fmt.Errorf("decode from json: %w", err)
	// }

	// for _, result := range resBody.Results {
	// 	if !result.Success {
	// 		return fmt.Errorf("a cancellation failed")
	// 	}
	// }

	return nil
}

func (e *Coinbase) ListOpenOrders(ctx context.Context) ([]Order, error) {
	url := fmt.Sprintf("%s/api/v3/brokerage/orders/historical/batch", coinbaseBaseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating req: %w", err)
	}

	query := req.URL.Query()
	query.Add("order_status", "OPEN")

	req.URL.RawQuery = query.Encode()

	res, err := e.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	type listOrder struct {
		OrderID   string `json: "order_id"`
		ProductID string `json: "product_id"`
		Side      string `json: "side"`
		ClientID  string `json: "client_id"`
	}

	type listOrdersResponse struct {
		Orders []listOrder `json: "orders"`
	}

	var resData listOrdersResponse

	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		return nil, fmt.Errorf("decode response body: %w", err)
	}

	orders := make([]Order, 0, len(resData.Orders))

	for _, order := range resData.Orders {
		side, err := e.convertToSide(order.Side)
		if err != nil {
			return nil, fmt.Errorf("failed to convert side %s: %w", order.Side, err)
		}

		pair, err := e.convertToPair(order.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to convert product ID %s: %w", order.ProductID, err)
		}

		orders = append(orders, Order{
			ID:       order.OrderID,
			ClientID: order.ClientID,
			Side:     side,
			Pair:     pair,
		})
	}

	return orders, nil
}

func (e *Coinbase) GetBalance(ctx context.Context, asset trading.Asset) (int64, error) {
	type availableBalance struct {
		Value    string `json: value:`
		Currency string `json: currency:`
	}

	type account struct {
		UUID             string           `json: uuid:`
		Name             string           `json: name:`
		Currency         string           `json: currency:`
		AvailableBalance availableBalance `json: "available_balance":`
	}

	type accountsReponse struct {
		Accounts []account `json: "accounts"`
	}

	url := fmt.Sprintf("%s/api/v3/brokerage/accounts", coinbaseBaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}

	res, err := e.doRequest(req)
	if err != nil {
		return 0, fmt.Errorf("perform request: %w", err)
	}

	var accounts accountsReponse

	if err = json.NewDecoder(res.Body).Decode(&accounts); err != nil {
		return 0, fmt.Errorf("decode body: %w", err)
	}

	var assetAccount account

	assetValue, err := e.convertAssetvalue(asset)
	if err != nil {
		return 0, fmt.Errorf("convert asset: %w", err)
	}

	for _, a := range accounts.Accounts {
		if a.Currency == assetValue {
			assetAccount = a
			break
		}
	}

	amount, err := asset.UnitStr(assetAccount.AvailableBalance.Value)
	if err != nil {
		return 0, fmt.Errorf("convert account amount: %w", err)
	}

	return amount, nil
}
