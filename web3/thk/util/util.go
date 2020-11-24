package util

import (
	"fmt"
	"math/big"
	"strings"
	"web3.go/common"

	"time"
)

type GetAccountJson struct {
	Address string `json:"address"`
	ChainId string `json:"chainId"`
}

type GetBlockTxsJson struct {
	ChainId string `json:"chainId"`
	Height  string `json:"height"`
	Page    string `json:"page"`
	Size    string `json:"size"`
}

type Account struct {
	Addr            string   `json:"address"`
	Nonce           uint64   `json:"nonce"`
	Balance         *big.Int `json:"balance"`         // 系统基础货币 TUE，不为nil
	LocalCurrency   *big.Int `json:"localCurrency"`   // 本链第二货币（如果存在），可为nil
	StorageRoot     []byte   `json:"storageRoot"`     // 智能合约使用的存储，Trie(key: Hash, value: Hash)
	CodeHash        []byte   `json:"codeHash"`        // 合约代码的Hash
	LongStorageRoot []byte   `json:"longStorageRoot"` // 系统合约用来保存更灵活的数据结构, Trie(key: Hash, value: []byte)
	ErrMsg          string   `json:"errMsg,omitempty"`
}
type Transaction struct {
	ChainId      string   `json:"chainId"`
	FromChainId  string   `json:"fromChainId,omitempty"`
	ToChainId    string   `json:"toChainId,omitempty"`
	From         string   `json:"from"`
	To           string   `json:"to"`
	Nonce        string   `json:"nonce"`
	Value        string   `json:"value"`
	Sig          string   `json:"sig,omitempty"`
	Pub          string   `json:"pub,omitempty"`
	Input        string   `json:"input"`
	UseLocal     bool     `json:"useLocal"`
	Extra        string   `json:"extra"` // 目前用来存交易类型，不存在时为普通交易，否则会对应特殊操作
	ExpireHeight int64    `json:"expireHeight,omitempty"`
	Multipubs    []string `json:"multipubs"`
	Multisigs    []string `json:"multisigs"`
}

// 此处与rpcTx的Hash算法一致
func (tx Transaction) HashValue() ([]byte, error) {
	s, err := tx.hashSerialize()
	if err != nil {
		return nil, err
	}
	return common.Hash256(s), nil
}

func (tx Transaction) hashSerialize() (string, error) {
	var toAddr string
	var fromAddr string
	if common.HasHexPrefix(tx.To) {
		toAddr = tx.To[2:]
		toAddr = strings.ToLower(toAddr)
	}

	if common.HasHexPrefix(tx.From) {
		fromAddr = tx.From[2:]
		fromAddr = strings.ToLower(fromAddr)
	}
	var input string
	if common.HasHexPrefix(tx.Input) {
		input = tx.Input[2:]
		input = strings.ToLower(input)
	}
	u := "0"
	if tx.UseLocal {
		u = "1"
	}
	if common.HasHexPrefix(tx.Extra) {
		tx.Extra = tx.Extra[2:]
	}
	extra := strings.ToLower(tx.Extra)
	str := []string{tx.ChainId, fromAddr, toAddr, tx.Nonce, u, tx.Value, input, extra}
	return strings.Join(str, "-"), nil
}

type GetTxByHash struct {
	ChainId string `json:"chainId"`
	Hash    string `json:"hash"`
}

type GetBlockHeader struct {
	ChainId string `json:"chainId"`
	Height  string `json:"height"`
}

type PingJson struct {
	Address string `json:"address"`
}

type GetChainInfoJson struct {
	ChainIds []int `json:"chainIds"`
}

type GetStatsJson struct {
	ChainId string `json:"chainId"`
}

type GetTransactionsJson struct {
	ChainId     string `json:"chainId"`
	Address     string `json:"address"`
	StartHeight string `json:"startHeight"`
	EndHeight   string `json:"endHeight"`
}

type GetMultiStatsJson struct {
	ChainId string `json:"chainId"`
}

type GetCommitteeJson struct {
	ChainId string `json:"chainId"`
	Epoch   string `json:"epoch"`
}
type CompileContractJson struct {
	ChainId  string `json:"chainId"`
	Contract string `json:"contract"`
}

type BlockGetCallback func() (interface{}, error)

func BlockGet(callback BlockGetCallback) (interface{}, error) {
	times := 1
	for {
		res, err := callback()
		if err != nil {
			return nil, err
		}
		if res == nil {
			if times > 5 {
				return nil, fmt.Errorf("get timeout")
			}
			times++
			time.Sleep(2 * time.Second)
			continue
		} else {
			return res, nil
		}
	}
}
