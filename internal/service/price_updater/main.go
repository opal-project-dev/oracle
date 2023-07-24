package price_updater

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/opal-project-dev/oracle/internal/config"
	"github.com/opal-project-dev/oracle/internal/service/contract"
	"github.com/opal-project-dev/oracle/internal/service/types"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"math/big"
	"time"
)

func PriceUpdater(cfg config.Config, input types.Data, nonce uint64) *priceUpdater {
	cfg.Log().Info("Running price_updater")
	client := cfg.ChainsInfo().DestClient
	contract, err := contract.NewContract(input.InternalAddress, client)
	if err != nil {
		cfg.Log().WithError(errors.WithMessage(err, "Failed to create new destination client entity.")).Error()
		return nil
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		cfg.Log().WithError(errors.WithMessage(err, "Failed to get destination network chain id.")).Error()
		return nil
	}

	return &priceUpdater{
		nonce:     int64(nonce),
		wallet:    cfg.WalletInfo(),
		client:    client,
		priceData: input.PriceData,
		contract:  contract,
		chainId:   chainId,
		interval:  cfg.RunnerInterval(),
		gasLimit:  cfg.ChainsInfo().GasLimit,
		log: cfg.Log().
			WithField("from", input.From).
			WithField("source", input.Source).
			WithField("price_feed_address", input.InternalAddress),
	}
}

type priceUpdater struct {
	nonce     int64
	wallet    config.WalletInfo
	client    *ethclient.Client
	chainId   *big.Int
	priceData types.PriceData
	contract  *contract.Contract
	interval  time.Duration
	gasLimit  uint64
	log       *logan.Entry
}
