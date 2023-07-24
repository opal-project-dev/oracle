package chainlink

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/opal-project-dev/oracle/internal/amount"
	"github.com/opal-project-dev/oracle/internal/config"
	"github.com/opal-project-dev/oracle/internal/service/contract"
	"github.com/opal-project-dev/oracle/internal/service/types"
	"github.com/pkg/errors"
)

type Chainlink struct {
	Cfg             config.Config
	ExternalAddress common.Address
	InternalAddress common.Address
}

func (c *Chainlink) GetCurrentPrice(priceChan chan<- types.Data) error {
	c.Cfg.Log().Info("Running chainlink listener")
	chainLink, err := contract.NewChainlink(c.ExternalAddress, c.Cfg.ChainsInfo().SourceClient)
	if err != nil {
		c.Cfg.Log().WithError(errors.WithMessage(err, "Failed to create new Chainlink entity.")).Error()
		return err
	}
	data, err := chainLink.LatestRoundData(&bind.CallOpts{})
	if err != nil {
		c.Cfg.Log().WithError(errors.WithMessage(err, "Failed to get latest data from ChainLink.")).Error()
		return err
	}
	roundData, err := chainLink.GetRoundData(&bind.CallOpts{}, data.RoundId)
	if err != nil {
		c.Cfg.Log().WithError(errors.WithMessage(err, "Failed to get round data from ChainLink.")).Error()
		return err
	}

	decimals, err := chainLink.Decimals(&bind.CallOpts{})
	if err != nil {
		return errors.Wrap(err, "Failed to get decimals from the contract.")
	}

	result := types.Data{
		Source:          "chainlink",
		From:            c.ExternalAddress.String(),
		InternalAddress: c.InternalAddress,
		PriceData: types.PriceData{
			RoundId:         roundData.RoundId,
			Answer:          amount.NewFromIntWithPrecision(roundData.Answer, int(decimals)),
			StartedAt:       roundData.StartedAt,
			UpdatedAt:       roundData.UpdatedAt,
			AnsweredInRound: roundData.AnsweredInRound,
		},
	}
	priceChan <- result
	return nil
}
