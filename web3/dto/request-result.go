package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"web3.go/web3/complex/types"
	"web3.go/web3/constants"
)

type RequestResult struct {
	// ID      int         `json:"id"`
	// Version string      `json:"jsonrpc"`
	Result interface{} `json:"result"`
	Error  *Error      `json:"error,omitempty"`
	Data   string      `json:"data,omitempty"`
}

type SendTxResult struct {
	TXhash string `json:"TXhash,omitempty"`
	ErrMsg string `json:"ErrMsg,omitempty"`
}
type RpcMakeVccProofJson struct {
	Proof  map[string]interface{} `json:"proof,omitempty"`
	ErrMsg string                 `json:"ErrMsg,omitempty"`
}

type MakeCCCExistenceProofJson struct {
	Proof  map[string]interface{} `json:"proof,omitempty"`
	ErrMsg string                 `json:"ErrMsg,omitempty"`
}

//GetCCCRelativeTx
type GetCCCRelativeTxJson struct {
	Proof  map[string]interface{} `json:"proof,omitempty"`
	ErrMsg string                 `json:"ErrMsg,omitempty"`
}
type CompileContractJson struct {
	Test   map[string]interface{} `json:"test,omitempty"`
	ErrMsg string                 `json:"ErrMsg,omitempty"`
}

type TransactionResult struct {
	ChainId   int      `json:"chainid"`
	From      string   `json:"from"`
	To        string   `json:"to"`
	Nonce     int      `json:"nonce"`
	Value     *big.Int `json:"value"`
	Input     string   `json:"input"`
	Hash      string   `json:"hash"`
	UseLocal  bool     `json:"uselocal"`
	Extra     string   `json:"extra"` // It is currently used to save transaction types. If it does not exist, it is a normal transaction. Otherwise, it will correspond to special operations
	Timestamp uint64   `json:"timestamp"`
}

type TxResult struct {
	Transaction     TransactionResult `json:"tx"`
	Root            string            `json:"root"`
	Status          int               `json:"status"`
	Logs            interface{}       `json:"logs"`
	TransactionHash string            `json:"transactionHash"`
	ContractAddress string            `json:"contractAddress"`
	Out             string            `json:"out"`
	GasFee          string            `json:"gasFee"`
	GasUsed         int               `json:"gasUsed"`
	BlockHeight     int               `json:"blockHeight"`
	Error           string            `json:"errorMsg"`
	ErrMsg          string            `json:"ErrMsg,omitempty"`
}

type GetBlockResult struct {
	Hash          string `json:"hash"`          // Hash of this block
	Previoushash  string `json:"previoushash"`  // Hash of parent block
	ChainId       int    `json:"chainid"`       //
	Height        int    `json:"height"`        // The block height of the query block
	Empty         bool   `json:"empty"`         // Is it an empty block
	RewardAddress string `json:"rewardaddress"` // Receiving address
	Mergeroot     string `json:"mergeroot"`     // Merge other chain block data hash
	Deltaroot     string `json:"deltaroot"`     // Cross chain transfer data hash
	Stateroot     string `json:"stateroot"`     // State hash
	RREra         int    `json:"rrera"`
	RRCurrent     string `json:"rrcurrent"`
	RRNext        string `json:"rrnext"`
	Txcount       int    `json:"txcount"`
	Timestamp     int64  `json:"timestamp"`
	ErrMsg        string `json:"ErrMsg,Omitempty"`
}

type NodeInfo struct {
	NodeId        string      `json:"nodeId"`        // Node ID
	Version       string      `json:"version"`       // edition
	IsDataNode    bool        `json:"isDataNode"`    // Is it a data node
	DataNodeOf    int         `json:"dataNodeOf"`    // Data node
	LastMsgTime   int64       `json:"lastMsgTime"`   // Last message time
	LastEventTime int64       `json:"lastEventTime"` // Last event time
	LastBlockTime int64       `json:"lastBlockTime"` // Last block time
	Overflow      bool        `json:"overflow"`      // overflow
	LastBlocks    interface{} `json:"lastBlocks"`    // The last block
	OpTypes       interface{} `json:"opTypes"`       // type
	ErrMsg        string      `json:"ErrMsg,Omitempty"`
}

type DataNode struct {
	DataNodeId   string `json:"dataNodeId"`
	DataNodeIp   string `json:"dataNodeIp"`
	DataNodePort int    `json:"dataNodePort"`
}
type GetChainInfo struct {
	ChainId   int        `json:"chainId"`
	DataNodes []DataNode `json:"datanodes"`
	Mode      int        `json:"mode"`
	Parent    int        `json:"parent"`
}

type BlockTxs struct {
	Elections      interface{}         `json:"elections"`
	AccountChanges []TransactionResult `json:"accountchanges"`
	ErrMsg         string              `json:"errMsg,omitempty"`
}

type GetTransactions struct {
	ChainId   int    `json:"chainId"`
	From      string `json:"from"`
	To        string `json:"to"`
	Nonce     int    `json:"nonce"`
	Value     int    `json:"value"`
	Input     string `json:"input"`
	Hash      string `json:"hash"`
	Timestamp int64  `json:"timestamp"`
}

type GetChainStats struct {
	ChainId           int      `json:"chainId"`
	CurrentHeight     int      `json:"currentheight"`
	EpochDuration     int      `json:"epochduration"`
	EpochLength       int      `json:"epochlength"`
	GasLimit          int      `json:"gaslimit"`
	GasPrice          string   `json:"gasprice"`
	LastEpochDuration int      `json:"lastepochduration"`
	Lives             int      `json:"lives"`
	Tps               int      `json:"tps"`
	TpsLastEpoch      int      `json:"tpsLastEpoch"`
	N                 int      `json:"n"`
	TpsLastN          int      `json:"tpsLastN"`
	LastNduration     int      `json:"lastNduration"`
	TxCount           int      `json:"txcount"`
	AccountCount      int      `json:"accountcount"`
	CurrentComm       []string `json:"currentcomm"`
}

type GetMultiStatsResult struct {
	ErrMsg string `json:"ErrMsg,Omitempty"`
}

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (pointer *RequestResult) ToStringArray() ([]string, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}
	result := (pointer).Result.([]interface{})
	new := make([]string, len(result))
	for i, v := range result {
		new[i] = v.(string)
	}
	return new, nil
}

func (pointer *RequestResult) ToComplexString() (types.ComplexString, error) {
	if err := pointer.checkResponse(); err != nil {
		return "", err
	}
	result := (pointer).Result.(interface{})
	return types.ComplexString(result.(string)), nil
}

func (pointer *RequestResult) ToString() (string, error) {
	if err := pointer.checkResponse(); err != nil {
		return "", err
	}
	result := (pointer).Result.(interface{})
	return result.(string), nil
}

func (pointer *RequestResult) ToInt() (int64, error) {
	if err := pointer.checkResponse(); err != nil {
		return 0, err
	}
	result := (pointer).Result.(interface{})
	hex := result.(string)
	numericResult, err := strconv.ParseInt(hex, 16, 64)
	return numericResult, err
}

func (pointer *RequestResult) ToBigInt() (*big.Int, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}
	res := (pointer).Result.(interface{})
	ret, success := big.NewInt(0).SetString(res.(string)[2:], 16)
	if !success {
		return nil, errors.New(fmt.Sprintf("Failed to convert %s to BigInt", res.(string)))
	}
	return ret, nil
}

func (pointer *RequestResult) ToComplexIntResponse() (types.ComplexIntResponse, error) {
	if err := pointer.checkResponse(); err != nil {
		return types.ComplexIntResponse(0), err
	}
	result := (pointer).Result.(interface{})
	var hex string
	switch v := result.(type) {
	// Testrpc returns a float64
	case float64:
		hex = strconv.FormatFloat(v, 'E', 16, 64)
		break
	default:
		hex = result.(string)
	}
	cleaned := strings.TrimPrefix(hex, "0x")
	return types.ComplexIntResponse(cleaned), nil
}

func (pointer *RequestResult) ToBoolean() (bool, error) {
	if err := pointer.checkResponse(); err != nil {
		return false, err
	}
	result := (pointer).Result.(interface{})
	return result.(bool), nil
}

func (pointer *RequestResult) ToSignTransactionResponse() (*SignTransactionResponse, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}
	result := (pointer).Result.(map[string]interface{})
	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}
	signTransactionResponse := &SignTransactionResponse{}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}
	err = json.Unmarshal(marshal, signTransactionResponse)
	return signTransactionResponse, err
}

func (pointer *RequestResult) ToTransactionResponse() (*TransactionResponse, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}
	result := (pointer).Result.(map[string]interface{})
	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}
	transactionResponse := &TransactionResponse{}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}
	err = json.Unmarshal(marshal, transactionResponse)
	return transactionResponse, err
}

func (pointer *RequestResult) ToTransactionReceipt() (*TransactionReceipt, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}
	result := (pointer).Result.(map[string]interface{})
	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}
	transactionReceipt := &TransactionReceipt{}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}
	err = json.Unmarshal(marshal, transactionReceipt)
	return transactionReceipt, err
}

func (pointer *RequestResult) ToBlock() (*Block, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}
	result := (pointer).Result.(map[string]interface{})
	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}
	block := &Block{}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}
	err = json.Unmarshal(marshal, block)
	return block, err
}

func (pointer *RequestResult) ToSyncingResponse() (*SyncingResponse, error) {
	if err := pointer.checkResponse(); err != nil {
		return nil, err
	}
	var result map[string]interface{}
	switch (pointer).Result.(type) {
	case bool:
		return &SyncingResponse{}, nil
	case map[string]interface{}:
		result = (pointer).Result.(map[string]interface{})
	default:
		return nil, customerror.UNPARSEABLEINTERFACE
	}
	if len(result) == 0 {
		return nil, customerror.EMPTYRESPONSE
	}
	syncingResponse := &SyncingResponse{}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nil, customerror.UNPARSEABLEINTERFACE
	}
	json.Unmarshal(marshal, syncingResponse)
	return syncingResponse, nil
}

// To avoid a conversion of a nil interface
func (pointer *RequestResult) checkResponse() error {
	if pointer.Error != nil {
		return errors.New(pointer.Error.Message)
	}
	if pointer.Result == nil {
		return customerror.EMPTYRESPONSE
	}
	return nil
}
