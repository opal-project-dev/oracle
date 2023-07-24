package types

import "github.com/pkg/errors"

type PriceSource string

const (
	Chainlink     = "chainlink"
	Coinmarketcap = "coinmarketcap"
	Coingecko     = "coingecko"
)

func (ps PriceSource) IsValid() error {
	switch ps {
	case Chainlink, Coinmarketcap, Coingecko:
		return nil
	}
	return errors.New("Invalid price source")
}
