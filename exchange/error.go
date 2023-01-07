package exchange

import "errors"

var (
	// ErrAPIKeyNotSet describes an error in which the API key is not set for
	// a client.
	ErrAPIKeyNotSet = errors.New("API Key not set")

	// ErrAPISecretNotSet describes an error in which the API secret is not
	// set for a client.
	ErrAPISecretNotSet = errors.New("API Secret not set")

	// ErrMissingPair describes an error that occurs when a pair has not
	// been implemented for an exchange.
	ErrMissingPair = errors.New("pair value is missing for exchange")
)
