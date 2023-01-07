package trading

import (
	"fmt"
	"math"
	"strconv"
)

// Asset represents a coin or currency that is used in trading pairs.
type Asset string

const (
	// BTC represents the Bitcoin asset
	BTC Asset = "BTC"
	// ETH represents the ethereum asset
	ETH Asset = "ETH"
	// USD represents the United States Dollar asset
	USD Asset = "USD"
)

// Decimals stores the number of decimal places that an asset has.
func (a Asset) Decimals() int {
	switch a {
	case BTC:
		const BTCDecimals = 8
		return BTCDecimals
	case ETH:
		const ETHDecimals = 18
		return ETHDecimals
	case USD:
		const USDDecimals = 2
		return USDDecimals
	}

	return 0
}

const floatBits = 64

// Unit returns the normalized unit value for an asset. i.e. for USD
// if the given value is: 50.00, then the normalized value is 5000, due
// to the number of cents in the value.
func (a Asset) Unit(u float64) int64 {
	return int64(u * math.Pow10(a.Decimals()))
}

// Format will convert the units into a string that has a floating point
// denomination.
func (a Asset) Format(i int64) string {
	const precision = -1

	f := float64(i) * math.Pow10(a.Decimals()*-1)

	return strconv.FormatFloat(f, 'f', precision, floatBits)
}

// UnitStr will produce a normalized unit from a string value. See
// Unit for more information.
func (a Asset) UnitStr(s string) (int64, error) {
	f, err := strconv.ParseFloat(s, floatBits)
	if err != nil {
		return 0, fmt.Errorf("parsing string: %w", err)
	}

	return a.Unit(f), nil
}
