package listener

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/opal-project-dev/oracle/internal/config"
	"github.com/opal-project-dev/oracle/internal/service/chainlink"
	"github.com/opal-project-dev/oracle/internal/service/coingecko"
	"github.com/opal-project-dev/oracle/internal/service/coinmarketcap"
	"github.com/opal-project-dev/oracle/internal/service/contract"
	"github.com/opal-project-dev/oracle/internal/service/source"
	"github.com/opal-project-dev/oracle/internal/service/types"
	"github.com/pkg/errors"
	"log"
	"math/big"
)

func NewListener(cfg config.Config, p config.Pairs) source.Source {
	switch p.Source {
	case types.Chainlink:
		return &chainlink.Chainlink{
			Cfg:             cfg,
			ExternalAddress: p.ExternalAddress,
			InternalAddress: p.InternalAddress,
		}
	case types.Coinmarketcap:
		return &coinmarketcap.Coinmarketcap{
			Cfg:                cfg,
			ApiPrivateKey:      p.ApiKey,
			CurrencyId:         p.CurrencyId,
			ConversionCurrency: p.ConversionCurrency,
			InternalAddress:    p.InternalAddress,
			RoundID:            getCurrentRound(cfg.ChainsInfo().DestClient, p.ExternalAddress),
		}
	case types.Coingecko:
		return &coingecko.Coingecko{
			Cfg:                cfg,
			CurrencyId:         p.CurrencyId,
			ConversionCurrency: p.ConversionCurrency,
			InternalAddress:    p.InternalAddress,
			RoundID:            getCurrentRound(cfg.ChainsInfo().DestClient, p.ExternalAddress),
		}
	default:
		panic("unknown source")
	}
}

func getCurrentRound(client *ethclient.Client, address common.Address) *big.Int {
	contr, err := contract.NewContract(address, client)
	if err != nil {
		log.Panic(errors.Wrap(err, "failed to create new contract entity"))
	}

	roundData, err := contr.LatestRoundData(&bind.CallOpts{})
	if err != nil {
		return big.NewInt(1)
	}

	return roundData.RoundId
}
