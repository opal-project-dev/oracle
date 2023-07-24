package config

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/opal-project-dev/oracle/internal/service/types"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"reflect"
	"strings"
)

type Pairer interface {
	Pairs() []Pairs
}

type Pairs struct {
	Source             types.PriceSource `fig:"source"`
	CurrencyId         string            `fig:"currency_id"`
	ConversionCurrency string            `fig:"conversion_currency"`
	ApiKey             string            `fig:"api_private_key"`
	ExternalAddress    common.Address    `fig:"external_address"`
	InternalAddress    common.Address    `fig:"internal_address,required"`
}

func NewPairs(getter kv.Getter) Pairer {
	return &pairer{
		getter: getter,
	}
}

type pairer struct {
	getter kv.Getter
	once   comfig.Once
}

func (p *pairer) Pairs() []Pairs {
	var config struct {
		P []Pairs `fig:"data"`
	}
	p.once.Do(func() interface{} {
		err := figure.
			Out(&config).
			With(figure.BaseHooks, AddressHook).
			From(kv.MustGetStringMap(p.getter, "pairs")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out contracts pairs"))
		}
		return config.P
	})
	return config.P
}

var AddressHook = figure.Hooks{
	"[]config.Pairs": func(value interface{}) (reflect.Value, error) {
		pairs := make([]Pairs, 0)
		switch s := value.(type) {
		case []interface{}:
			var err error
			for _, elem := range s {
				mapa, ok := elem.(map[interface{}]interface{})
				if !ok {
					return reflect.Value{}, errors.Wrap(err,
						"failed to cast mapa to interface")
				}

				var data Pairs
				contractAddress, ok := mapa["internal_address"].(string)
				if !ok {
					return reflect.Value{}, errors.Wrap(err, "failed to get internal smart contract address")
				}

				data.InternalAddress = common.HexToAddress(contractAddress)

				source := types.PriceSource(strings.ToLower(mapa["source"].(string)))
				if err := source.IsValid(); err != nil {
					return reflect.Value{}, errors.Wrap(err, "failed to get price source.")
				}
				data.Source = source

				switch source {
				case types.Chainlink:
					contractAddress, ok = mapa["external_address"].(string)
					if !ok {
						return reflect.Value{}, errors.Wrap(err, "failed to get chainlink smart contract address")
					}
					data.ExternalAddress = common.HexToAddress(contractAddress)
				case types.Coinmarketcap:
					apiKey, ok := mapa["api_private_key"].(string)
					if !ok {
						return reflect.Value{}, errors.Wrap(err, "failed to get api private key")
					}
					data.ApiKey = apiKey
					id, ok := mapa["currency_id"].(string)
					if !ok {
						return reflect.Value{}, errors.Wrap(err, "failed to get currency id")
					}
					data.CurrencyId = id
					conversionCurrency, ok := mapa["conversion_currency"].(string)
					if !ok {
						return reflect.Value{}, errors.Wrap(err, "failed to get conversion currency")
					}
					data.ConversionCurrency = conversionCurrency
				case types.Coingecko:
					id, ok := mapa["currency_id"].(string)
					if !ok {
						return reflect.Value{}, errors.Wrap(err, "failed to get pairs string")
					}
					data.CurrencyId = id
					conversionCurrency, ok := mapa["conversion_currency"].(string)
					if !ok {
						return reflect.Value{}, errors.Wrap(err, "failed to get pairs string")
					}
					data.ConversionCurrency = conversionCurrency

				default:
					return reflect.Value{}, errors.Errorf("there is no such source: %s", source)
				}

				pairs = append(pairs, data)
			}
			return reflect.ValueOf(pairs), nil
		default:
			return reflect.ValueOf(value), errors.New("failed to get case")
		}
	},
}
