package config

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"reflect"
)

type ChainsInitiator interface {
	ChainsInfo() ChainsInfo
}

func NewInitiator(getter kv.Getter) ChainsInitiator {
	return &connector{
		getter: getter,
	}
}

type connector struct {
	getter kv.Getter
	once   comfig.Once
}

type ChainsInfo struct {
	SourceClient *ethclient.Client `fig:"source_network,required"`
	DestClient   *ethclient.Client `fig:"destination_network,required"`
	GasLimit     uint64            `fig:"gas_limit,required"`
}

func (c *connector) ChainsInfo() ChainsInfo {
	return c.once.Do(func() interface{} {
		cfg := ChainsInfo{}

		err := figure.
			Out(&cfg).
			With(figure.BaseHooks, ChainsHook).
			From(kv.MustGetStringMap(c.getter, "chains")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out chains"))
		}
		return cfg
	}).(ChainsInfo)
}

var ChainsHook = figure.Hooks{
	"*ethclient.Client": func(value interface{}) (reflect.Value, error) {
		endpoint, err := cast.ToStringE(value)
		if err != nil {
			return reflect.Value{}, err
		}

		dial, err := ethclient.Dial(endpoint)
		if err != nil {
			panic(errors.Wrap(err, "failed to dial ethclient"))
		}
		return reflect.ValueOf(dial), nil
	},
}
