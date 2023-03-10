package exchange_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/project-code-io/crypto-trading-bot-go/exchange"
)

func TestBinanceConstructor(t *testing.T) {
	testCases := []struct {
		name  string
		input exchange.BinanceDomain
		wants exchange.Binance
	}{
		{
			name:  "testing that the BinanceDomainUS has the correct baseURL",
			input: exchange.BinanceDomainUS,
			wants: exchange.Binance{
				BaseURL: "https://api.binance.us",
			},
		},
		{
			name:  "testing that the BinanceDomainDotCom has the correct baseURL",
			input: exchange.BinanceDomainDotCom,
			wants: exchange.Binance{
				BaseURL: "https://api.binance.com",
			},
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			res := exchange.NewBinance(tt.input)

			assert.Equal(t, tt.wants, *res)
		})
	}
}
