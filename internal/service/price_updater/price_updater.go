package price_updater

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"math/big"
)

func (i *priceUpdater) Update() (txSent bool, err error) {
	if i == nil {
		return false, errors.New("Failed to create price updater entity")
	}

	gasPrice, err := i.client.SuggestGasPrice(context.Background())
	if err != nil {
		i.log.WithError(errors.WithMessage(err, "Failed to get destination network gas price.")).Error()
		return false, err
	}
	roundId := i.priceData.RoundId
	if i.priceData.RoundId == nil {
		latestRoundData, err := i.contract.LatestRoundData(&bind.CallOpts{})
		if err != nil {
			i.log.WithError(errors.WithMessage(err, "Failed to get latest round data from the contract.")).Error()
			return false, err
		}
		if latestRoundData.RoundId.Cmp(big.NewInt(0)) == 0 {
			roundId = big.NewInt(0)
		}
		roundId = latestRoundData.RoundId.Add(latestRoundData.RoundId, big.NewInt(1))
	}

	decimals, err := i.contract.Decimals(&bind.CallOpts{})
	if err != nil {
		return false, errors.Wrap(err, "Failed to get decimals from the contract.")
	}

	amount := i.priceData.Answer.IntWithPrecision(int(decimals))

	data, err := i.contract.LatestRoundData(&bind.CallOpts{})
	if err != nil {
		return false, err
	}
	roundData, err := i.contract.GetRoundData(&bind.CallOpts{}, data.RoundId)
	if err != nil {
		return false, err
	}

	if amount.Cmp(roundData.Answer) == 0 {
		i.log.WithField("price", amount).Info("Price is up to date.")
		return false, nil
	}

	tx, err := i.contract.Set(i.txOpts(gasPrice, big.NewInt(i.nonce)), roundId, amount)
	if err != nil {
		i.log.WithError(errors.WithMessage(err, "Failed to set price to destination network contract.")).Error()
		return false, err
	}
	i.log.Infof("Successful transaction %s", tx.Hash().Hex())
	return true, err
}

func (i *priceUpdater) txOpts(gasPrice, nonce *big.Int) *bind.TransactOpts {
	return &bind.TransactOpts{
		From:  i.wallet.Address,
		Nonce: nonce,
		Signer: func(address common.Address, tx *ethtypes.Transaction) (*ethtypes.Transaction, error) {
			signer := ethtypes.NewEIP155Signer(i.chainId)
			signature, err := crypto.Sign(signer.Hash(tx).Bytes(), i.wallet.PrivateKey)
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
		Value:     nil,
		GasPrice:  gasPrice,
		GasFeeCap: nil,
		GasTipCap: nil,
		GasLimit:  i.gasLimit,
		Context:   nil,
		NoSend:    false,
	}
}
