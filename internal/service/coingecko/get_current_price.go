package coingecko

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/opal-project-dev/oracle/internal/config"
	"github.com/opal-project-dev/oracle/internal/service/types"
	"math/big"
	"time"
)

type Coingecko struct {
	Cfg                config.Config
	CurrencyId         string
	ConversionCurrency string
	InternalAddress    common.Address
	RoundID            *big.Int
}

func (c *Coingecko) GetCurrentPrice(priceChan chan<- types.Data) error {
	c.Cfg.Log().Info("Running coingecko listener")
	price, err := getCurrentPrice(c.CurrencyId, c.ConversionCurrency)
	if err != nil {
		c.Cfg.Log().WithError(err).Error()
		return err
	}
	var data types.Data
	data.PriceData.Answer = price
	c.RoundID.Add(c.RoundID, big.NewInt(1))
	data.PriceData.RoundId = c.RoundID
	data.PriceData.StartedAt = big.NewInt(time.Now().Unix())
	data.PriceData.UpdatedAt = big.NewInt(time.Now().Unix())
	data.PriceData.AnsweredInRound = c.RoundID
	result := types.Data{
		Source:          "coingecko",
		From:            c.CurrencyId,
		InternalAddress: c.InternalAddress,
		PriceData:       data.PriceData,
	}
	priceChan <- result
	return nil
}
