package price_updater

import (
	"github.com/opal-project-dev/oracle/internal/config"
	"github.com/opal-project-dev/oracle/internal/service/types"
	"strings"
	"time"
)

func NewPriceUpdater(cfg config.Config, input chan types.Data, nonce uint64) error {
	cfg.Log().Info("Running price_updater")
	for x := range input {
		txSent, err := PriceUpdater(cfg, x, nonce).Update()
		if err != nil {
			if strings.Contains(err.Error(), "Pool(TemporarilyBanned)") {
				cfg.Log().Warn("Banned")
				time.Sleep(1 * time.Minute)
				continue
			}
			if strings.Contains(err.Error(), "Pool(InvalidTransaction(InvalidTransaction::Stale))") ||
				strings.Contains(err.Error(), "Pool(TooLowPriority") {
				cfg.Log().WithField("nonce", nonce).Error("Nonce is too low")
				nonce++
				continue
			}
			cfg.Log().Warn("Contract adress:", x.InternalAddress, " Nonce:", nonce)
			cfg.Log().WithError(err).Error("Failed to update price")
			continue
		}
		if txSent {
			nonce++
		}
	}
	return nil
}
