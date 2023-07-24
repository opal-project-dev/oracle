package coinmarketcap

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/opal-project-dev/oracle/internal/config"
	"github.com/opal-project-dev/oracle/internal/service/types"
	"math/big"
	"time"
)

type Coinmarketcap struct {
	Cfg                config.Config
	ApiPrivateKey      string
	CurrencyId         string
	ConversionCurrency string
	InternalAddress    common.Address
	RoundID            *big.Int
}

func (c *Coinmarketcap) GetCurrentPrice(priceChan chan<- types.Data) error {
	c.Cfg.Log().Info("Running coinmarketcap listener")
	price, err := getCurrentPrice(c.CurrencyId, c.ConversionCurrency, c.ApiPrivateKey)
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
		Source:          "coinmarketcap",
		From:            c.CurrencyId,
		InternalAddress: c.InternalAddress,
		PriceData:       data.PriceData,
	}
	priceChan <- result
	return nil

}
