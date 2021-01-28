package test

import (
	"errors"
	"fmt"
	"web3.go/web3/dto"
	"web3.go/web3/thk/util"
)

func BlockGetTransactionReceipt(chainId, hash string) (*dto.TxResult, error) {
	receipt := util.BlockGetDefault(func() (recipt interface{}, _break bool) {
		recipt, err := Web3.Thk.GetTransactionByHash(chainId, hash)
		if err == nil {
			return recipt, true
		}
		return nil, false
	})
	if receipt == nil {
		return nil, fmt.Errorf("get transaction receipt timeout, chainId:%s, hash:%s", chainId, hash)
	}
	return receipt.(*dto.TxResult), nil
}

func BlockCheckTransactionReceipt(chainId, hash string) error {
	receipt, err := BlockGetTransactionReceipt(chainId, hash)
	if err != nil {
		return err
	}
	if receipt.Status != 1 {
		return errors.New(fmt.Sprintf("tx failed, hash:%v", hash))
	}
	return nil
}
