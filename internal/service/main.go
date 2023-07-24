package service

import (
	"context"
	"github.com/opal-project-dev/oracle/internal/config"
	"github.com/opal-project-dev/oracle/internal/service/listener"
	"github.com/opal-project-dev/oracle/internal/service/price_updater"
	"github.com/opal-project-dev/oracle/internal/service/types"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	"sync"
)

type service struct {
	cfg config.Config
	log *logan.Entry
}

func (s *service) run() {
	s.log.Info("Running service")
	ctx := context.Background()
	var wg sync.WaitGroup
	client := s.cfg.ChainsInfo().DestClient
	nonce, err := client.PendingNonceAt(context.Background(), s.cfg.WalletInfo().Address)
	if err != nil {
		s.log.WithError(errors.WithMessage(err, "failed to get pending nonce.")).Error()
		return
	}
	priceChan := make(chan types.Data)
	for _, p := range s.cfg.Pairs() {
		wg.Add(1)
		p := p
		go func() {
			defer wg.Done()
			running.WithBackOff(ctx, s.cfg.Log(), "new-listener-service", func(ctx context.Context) error {
				err = listener.NewListener(s.cfg, p).GetCurrentPrice(priceChan)
				return err
			}, s.cfg.RunnerInterval(), s.cfg.RunnerInterval(), s.cfg.RunnerInterval()*10)

		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = price_updater.NewPriceUpdater(s.cfg, priceChan, nonce)
		if err != nil {
			s.cfg.Log().WithError(errors.WithMessage(err, "Failed to update price.")).Error()
			return
		}
	}()
	wg.Wait()
}

func newService(cfg config.Config) *service {
	return &service{
		log: cfg.Log(),
		cfg: cfg,
	}
}

func Run(cfg config.Config) {
	newService(cfg).run()
}
