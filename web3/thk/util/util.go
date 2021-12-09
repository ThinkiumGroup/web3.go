package util

import (
	"encoding/json"
	"errors"
	common2 "github.com/ThinkiumGroup/go-common"
	"github.com/ThinkiumGroup/web3.go/common"
	"math/big"
	"strconv"
	"strings"

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
	Balance         *big.Int `json:"balance"`         // System base currency TKM, not nil
	LocalCurrency   *big.Int `json:"localCurrency"`   // The second currency in the chain (if any) can be nil
	StorageRoot     []byte   `json:"storageRoot"`     // Storage used by smart contract, trie (key: hash, value: hash)
	CodeHash        []byte   `json:"codeHash"`        // Hash of contract code
	LongStorageRoot []byte   `json:"longStorageRoot"` // System contract is used to save more flexible data structure, trie (key: hash, value: [] byte)
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
	Extra        string   `json:"extra"` // It is currently used to save transaction types. If it does not exist, it is a normal transaction. Otherwise, it will correspond to special operations
	ExpireHeight int64    `json:"expireHeight,omitempty"`
	Multipubs    []string `json:"multipubs"`
	Multisigs    []string `json:"multisigs"`
}

type GasProvider struct {
	Type     byte     `json:"type"`
	Gas      uint64   `json:"gas"`
	GasPrice *big.Int `json:"gasPrice"`
}

var DefaultGasProvider = &GasProvider{
	Type:     0,
	Gas:      25000,
	GasPrice: big.NewInt(40 * 10000 * 10000),
}

var BaseChainId int64 = 100007

func (tx *Transaction) HashValue() ([]byte, error) {
	chainId, ok := new(big.Int).SetString(tx.ChainId, 10)
	if !ok {
		return nil, errors.New("error chainId")
	}
	chainId.Add(chainId, big.NewInt(BaseChainId))
	if tx.Value == "" {
		tx.Value = "0"
	}
	value, ok := new(big.Int).SetString(tx.Value, 10)
	if !ok {
		return nil, errors.New("error value")
	}
	var gasProvider GasProvider
	gasBytes := common.FromHex(tx.Extra)
	if len(gasBytes) > 0 {
		err := json.Unmarshal(gasBytes, &gasProvider)
		if err != nil {
			return nil, err
		}
	}
	if gasProvider.GasPrice == nil {
		gasProvider.GasPrice = big.NewInt(0)
	}
	nonce, err := strconv.ParseInt(tx.Nonce, 10, 64)
	if err != nil {
		return nil, err
	}
	hash := common2.RlpHash([]interface{}{
		uint64(nonce),
		gasProvider.GasPrice,
		gasProvider.Gas,
		common.FromHex(tx.To),
		value,
		common.FromHex(tx.Input),
		chainId, uint(0), uint(0),
	})
	return hash.Bytes(), nil
}

// Deprecated
func (tx Transaction) hashSerialize() (string, error) {
	toAddr := strings.ToLower(common.CleanHexPrefix(tx.To))
	fromAddr := strings.ToLower(common.CleanHexPrefix(tx.From))
	input := strings.ToLower(common.CleanHexPrefix(tx.Input))
	u := "0"
	if tx.UseLocal {
		u = "1"
	}
	extra := strings.ToLower(common.CleanHexPrefix(tx.Extra))
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

type Callback func() (res interface{}, _break bool)

func BlockGetDefault(callback Callback) interface{} {
	return BlockGet(callback, 10, 3)
}
func BlockGet(callback Callback, maxTime, sleepSeconds int) interface{} {
	times := 1
	for {
		res, _break := callback()
		if !_break && times < maxTime && res == nil {
			times++
			time.Sleep(time.Duration(sleepSeconds) * time.Second)
			continue
		}
		return res
	}
}
