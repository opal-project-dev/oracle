package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/opal-project-dev/oracle/internal/amount"
	"math/big"
)

type Data struct {
	Source          string
	From            string
	InternalAddress common.Address
	PriceData       PriceData
}

type PriceData struct {
	RoundId         *big.Int
	Answer          amount.Amount
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}
