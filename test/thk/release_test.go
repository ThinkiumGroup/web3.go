package test

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"strconv"
	"testing"
	"time"
	"web3.go/common"
	"web3.go/common/hexutil"
	"web3.go/test"
	"web3.go/web3/thk/util"
)

func TestRelease(t *testing.T) {
	from := "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"
	contentTokenVesting, err := ioutil.ReadFile(test.TokenVestingJsonLocation)
	if err != nil {
		println(err)

	}
	var TokenVestingResponse ContractAbi
	_ = json.Unmarshal(contentTokenVesting, &TokenVestingResponse)

	//bytecodeTokenVesting := TokenVestingResponse.Bytecode
	contractTokenVesting, err := test.Web3.Thk.NewContract(TokenVestingResponse.Abi)
	hashTokenVesting := "0xe0c77573eb6e3447dcd7c361a9c6ef91d51e5376520664d2f6aad3aed04a2ae6"
	receiptTokenVesting, err := test.Web3.Thk.GetTransactionByHash(test.Web3.Thk.DefaultChainId, hashTokenVesting)
	toTokenVesting := receiptTokenVesting.ContractAddress

	//Vesting
	nonce, err := test.Web3.Thk.GetNonce(from, test.Web3.Thk.DefaultChainId)
	if err != nil {
		println("get nonce error")
		return
	}
	transactionTokenVesting := util.Transaction{
		ChainId: test.Web3.Thk.DefaultChainId, FromChainId: test.Web3.Thk.DefaultChainId, ToChainId: test.Web3.Thk.DefaultChainId, From: from,
		To: toTokenVesting, Value: "0", Input: "", Nonce: strconv.Itoa(int(nonce)),
	}
	tmcliff, errc := time.Parse("2006-01-02 15:04:05", "2019-07-19 17:47:00")
	tmstart, errc := time.Parse("2006-01-02 15:04:05", "2019-07-19 17:48:00")
	tmend, errc := time.Parse("2006-01-02 15:04:05", "2019-07-19 17:49:00")
	println(errc)
	toAddress, err := hexutil.Decode("0x14723a09acff6d2a60dcdf7aa4aff308fddc160c")
	vestAddress := common.BytesToAddress(toAddress)
	//vestAddressP,boolerr := new(big.Int).SetString(strings.TrimPrefix(string("0x14723a09acff6d2a60dcdf7aa4aff308fddc160c"), "0x"),16)
	//println(boolerr)
	cliff := tmcliff.Unix()
	start := tmstart.Unix()
	end := tmend.Unix()
	tmcliffP := new(big.Int).SetInt64(cliff)
	tmstartP := new(big.Int).SetInt64(start)
	tmendP := new(big.Int).SetInt64(end)
	timesP := new(big.Int).SetInt64(2)
	total := new(big.Int).SetUint64(uint64(100))
	resultTokenVesting, err := contractTokenVesting.Send(transactionTokenVesting, "addPlan", test.TmpKey,
		vestAddress, tmcliffP, tmstartP, timesP, tmendP, total, false, "上交易所后私募锁仓10分钟，之后每10分钟释放50%")
	if err != nil {
		t.Error(err)
	}
	t.Log("result:", resultTokenVesting)

	nonce, err = test.Web3.Thk.GetNonce(from, test.Web3.Thk.DefaultChainId)
	if err != nil {
		println(err)
		println("get nonce error")
		return
	}
	transactionReleasable := util.Transaction{
		ChainId: test.Web3.Thk.DefaultChainId, FromChainId: test.Web3.Thk.DefaultChainId, ToChainId: test.Web3.Thk.DefaultChainId, From: from,
		To: toTokenVesting, Value: "0", Input: "", Nonce: strconv.Itoa(int(nonce)),
	}
	value := new(big.Int).SetUint64(uint64(0))
	result, err := contractTokenVesting.Call(transactionReleasable, "releasableAmount", test.TmpKey, vestAddress)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("result:", result)
	err = contractTokenVesting.Parse(result.Out, "releasableAmount", &value)
	if err != nil {
		println("failed")
		return
	}

	transactionRelease := util.Transaction{
		ChainId: test.Web3.Thk.DefaultChainId, FromChainId: test.Web3.Thk.DefaultChainId, ToChainId: test.Web3.Thk.DefaultChainId, From: from,
		To: toTokenVesting, Value: "0", Input: "", Nonce: strconv.Itoa(int(nonce)),
	}

	hashTokenVesting, err = contractTokenVesting.Send(transactionRelease, "release", test.TmpKey, vestAddress)
	if err != nil {
		t.Error(err)
	}
	t.Log("result:", hashTokenVesting)
}
