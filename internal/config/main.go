package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
)

type Config interface {
	ChainsInitiator
	Runner
	comfig.Logger
	types.Copuser
	comfig.Listenerer
	Walleter
	Pairer
}

type config struct {
	ChainsInitiator
	comfig.Logger
	types.Copuser
	comfig.Listenerer
	getter kv.Getter
	Walleter
	Runner
	Pairer
}

func New(getter kv.Getter) Config {
	return &config{
		getter:          getter,
		Copuser:         copus.NewCopuser(getter),
		Listenerer:      comfig.NewListenerer(getter),
		Logger:          comfig.NewLogger(getter, comfig.LoggerOpts{}),
		Walleter:        NewWalletInfo(getter),
		Runner:          NewRunnerInfo(getter),
		ChainsInitiator: NewInitiator(getter),
		Pairer:          NewPairs(getter),
	}
}
