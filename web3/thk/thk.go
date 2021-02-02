package thk

import (
	"errors"
	"fmt"
	"math/big"
	"web3.go/common"
	"web3.go/common/hexutil"
	"web3.go/web3/dto"
	"web3.go/web3/providers"
	"web3.go/web3/thk/util"
)

type Thk struct {
	DefaultAddress          string
	DefaultPrivateKey       string
	DefaultExtraPrivateKeys []string
	DefaultAuthKey          string
	DefaultChainId          string

	provider providers.ProviderInterface
}

func NewThk(provider providers.ProviderInterface) *Thk {
	thk := new(Thk)
	thk.provider = provider
	return thk
}

func (thk *Thk) GetAccount(address string, chainId string) (*util.Account, error) {
	params := util.GetAccountJson{
		Address: address,
		ChainId: chainId,
	}
	res := util.Account{}
	if err := thk.provider.SendRequest(&res, "GetAccount", params); err != nil {
		return nil, err
	}
	if res.ErrMsg != "" {
		return nil, errors.New(res.ErrMsg)
	}
	return &res, nil
}

func (thk *Thk) GetBalance(address string, chainId string) (*big.Int, error) {
	res, err := thk.GetAccount(address, chainId)
	if err != nil {
		return nil, err
	}
	ret := res.Balance
	return ret, nil
}

func (thk *Thk) GetNonce(address string, chainId string) (int64, error) {
	res, err := thk.GetAccount(address, chainId)
	if err != nil {
		return 0, err
	}
	return int64(res.Nonce), nil
}

func (thk *Thk) GetBlockTxs(chainId string, height string, page string, size string) (*dto.BlockTxs, error) {
	params := util.GetBlockTxsJson{
		ChainId: chainId,
		Height:  height,
		Page:    page,
		Size:    size,
	}
	var res dto.BlockTxs
	if err := thk.provider.SendRequest(&res, "GetBlockTxs", params); err != nil {
		return nil, err
	}
	if res.ErrMsg != "" {
		return nil, errors.New(res.ErrMsg)
	}
	return &res, nil
}

func (thk *Thk) SendTx(transaction *util.Transaction) (string, error) {
	res := new(dto.SendTxResult)
	if err := thk.provider.SendRequest(res, "SendTx", transaction); err != nil {
		return "", err
	}
	if res.ErrMsg != "" {
		err := errors.New(res.ErrMsg)
		return "", err
	}
	return res.TXhash, nil
}

func (thk *Thk) SignTransaction(transaction *util.Transaction, privateKey string, multikeys ...string) error {
	hash, err := transaction.HashValue()
	if err != nil {
		return err
	}
	key, err := common.HexToPrivateKey(privateKey)
	if err != nil {
		return err
	}

	sig, err := common.Cipher.Sign(common.Cipher.PrivToBytes(key), hash)
	if err != nil {
		return err
	}

	transaction.Sig = hexutil.Encode(sig)
	transaction.Pub = hexutil.Encode(key.GetPublicKey().ToBytes())

	if len(multikeys) > 0 {
		for i := 0; i < len(multikeys); i++ {
			key, err = common.HexToPrivateKey(multikeys[i])
			if err != nil {
				return err
			}
			sign, err := common.Cipher.Sign(common.Cipher.PrivToBytes(key), hash)
			if err != nil {
				return err
			}
			transaction.Multisigs = append(transaction.Multisigs, hexutil.Encode(sign))
			transaction.Multipubs = append(transaction.Multipubs, hexutil.Encode(key.GetPublicKey().ToBytes()))
		}
	}
	return nil
}

func (thk *Thk) CallTransaction(transaction *util.Transaction) (*dto.TxResult, error) {
	res := new(dto.TxResult)
	if err := thk.provider.SendRequest(res, "CallTransaction", transaction); err != nil {
		return nil, err
	}
	if res.ErrMsg != "" {
		err := errors.New(res.ErrMsg)
		return nil, err
	}
	return res, nil
}

func (thk *Thk) GetTransactionByHash(chainId string, hash string) (*dto.TxResult, error) {
	params := util.GetTxByHash{
		ChainId: chainId,
		Hash:    hash,
	}
	res := new(dto.TxResult)
	if err := thk.provider.SendRequest(res, "GetTransactionByHash", params); err != nil {
		return nil, err
	}
	if res.ErrMsg != "" {
		err := errors.New(res.ErrMsg)
		return nil, err
	}
	return res, nil
}

func (thk *Thk) GetBlockHeader(chainId string, height string) (*dto.GetBlockResult, error) {
	params := util.GetBlockHeader{
		ChainId: chainId,
		Height:  height,
	}
	res := new(dto.GetBlockResult)
	if err := thk.provider.SendRequest(res, "GetBlockHeader", params); err != nil {
		return nil, err
	}
	if res.ErrMsg != "" {
		err := errors.New(res.ErrMsg)
		return nil, err
	}
	return res, nil
}

// Ping
func (thk *Thk) Ping(address string) (*dto.NodeInfo, error) {
	params := util.PingJson{
		Address: address,
	}
	res := new(dto.NodeInfo)
	if err := thk.provider.SendRequest(&res, "/chaininfo:Ping", params); err != nil {
		return nil, err
	}

	if res.ErrMsg != "" {
		return nil, errors.New(res.ErrMsg)
	}
	return res, nil
}

func (thk *Thk) GetChainInfo(chainIds []int) ([]dto.GetChainInfo, error) {
	params := new(util.GetChainInfoJson)
	params.ChainIds = chainIds
	var resArray []dto.GetChainInfo
	if err := thk.provider.SendRequest(&resArray, "/chaininfo:GetChainInfo", params); err != nil {
		return nil, err
	}
	return resArray, nil
}

func (thk *Thk) GetStats(chainId string) (gts dto.GetChainStats, err error) {
	params := new(util.GetStatsJson)
	params.ChainId = chainId
	res := new(dto.GetChainStats)
	if err := thk.provider.SendRequest(&res, "GetStats", params); err != nil {
		return *res, err
	}
	return *res, nil
}

// GetTransactions
func (thk *Thk) GetTransactions(chainId, address, startHeight, endHeight string) ([]dto.GetTransactions, error) {
	params := util.GetTransactionsJson{
		ChainId:     chainId,
		Address:     address,
		StartHeight: startHeight,
		EndHeight:   endHeight,
	}

	res := new(dto.GetTransactions)
	if err := thk.provider.SendRequest(res, "GetTransactions", params); err != nil {
		return nil, err
	}

	resArray := []dto.GetTransactions{*res}
	return resArray, nil
}

func (thk *Thk) GetCommittee(chainId string, epoch string) ([]string, error) {
	params := util.GetCommitteeJson{
		ChainId: chainId,
		Epoch:   epoch,
	}
	var res []string
	if err := thk.provider.SendRequest(&res, "/chaininfo:GetCommittee", params); err != nil {
		return nil, err
	}
	return res, nil
}

func (thk *Thk) RpcMakeVccProof(cashCheque *CashCheque) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	if err := thk.provider.SendRequest(&res, "RpcMakeVccProof", &cashCheque); err != nil {
		return nil, err
	}
	return res, nil
}

func (thk *Thk) MakeCCCExistenceProof(cashCheque *CashCheque) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	if err := thk.provider.SendRequest(&res, "MakeCCCExistenceProof", cashCheque); err != nil {
		return nil, err
	}
	if res["errMsg"] != nil && res["errMsg"].(string) != "" {
		err := errors.New(res["errMsg"].(string))
		return nil, err
	}
	return res, nil
}

// GetCCCRelativeTx
func (thk *Thk) GetCCCRelativeTx(transaction *util.Transaction) (map[string]interface{}, error) {
	res := new(dto.GetCCCRelativeTxJson)
	if err := thk.provider.SendRequest(res, "GetCCCRelativeTx", transaction); err != nil {
		return nil, err
	}
	if res.ErrMsg != "" {
		err := errors.New(res.ErrMsg)
		return nil, err
	}
	return res.Proof, nil
}

// get nodeSig  nodeId,  address bindAddr privateKey for hex with 0x
//  nodeType  should be 0 for Consensus, 1 for data
//  nonce  amount   string
func (thk *Thk) GetNodeSig(nodeId string, nodeType string, address string, nonce string, amount string, privateKey string) (string, error) {
	str := fmt.Sprintf("%s,%s,%s,%s,%s", nodeId[2:], nodeType, address[2:], nonce, amount)
	sign, err := common.Sign(str, privateKey)
	if err != nil {
		return "", err
	}
	return sign, nil
}
