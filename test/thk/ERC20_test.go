package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"
	"web3.go/common"
	"web3.go/common/hexutil"
	"web3.go/test"
	"web3.go/test/compiler"
	"web3.go/web3/thk/util"
)

func TestCompileERC20(t *testing.T) {
	res, err := compiler.CompileContract("../resources/contract/ERC20.sol",
		"../resources/contract/IERC20.sol", "../resources/contract/Pausable.sol",
		"../resources/contract/SafeMath.sol", "../resources/contract/Ownable.sol",
		"../resources/contract/Address.sol", "../resources/contract/SafeERC20.sol",
		"../resources/contract/TokenVesting.sol")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	var ERC20Json, TokenVestingJson ContractJson
	ERC20Json, TokenVestingJson, err = GetJson(res)
	dataERC20, err := json.MarshalIndent(ERC20Json, "", "  ")
	if ioutil.WriteFile(test.Erc20JsonLocation, dataERC20, 0644) == nil {
		fmt.Println("写入文件成功")
	}
	dataTokenVesting, err := json.MarshalIndent(TokenVestingJson, "", "  ")
	if ioutil.WriteFile(test.TokenVestingJsonLocation, dataTokenVesting, 0644) == nil {
		fmt.Println("写入文件成功")
	}
}
func GetJson(ctct map[string]interface{}) (ContractJson, ContractJson, error) {
	var contractJson, ERC20Json, TokenVestingJson ContractJson
	for keyname, value := range ctct {
		contractJson.ContractName = keyname
		arr := strings.Split(contractJson.ContractName, ":")
		length := len(arr) - 1
		if arr[length] == "ERC20" {
			mapvalue := value.(map[string]interface{})
			ERC20Json.ByteCode = mapvalue["code"].(string)
			info := mapvalue["info"].(map[string]interface{})
			abidef := info["abiDefinition"]
			abibytes, _ := json.Marshal(abidef)
			ERC20Json.ABI = string(abibytes)
		}
		if arr[length] == "TokenVesting" {
			mapvalue := value.(map[string]interface{})
			TokenVestingJson.ByteCode = mapvalue["code"].(string)
			info := mapvalue["info"].(map[string]interface{})
			abidef := info["abiDefinition"]
			abibytes, _ := json.Marshal(abidef)
			TokenVestingJson.ABI = string(abibytes)
		}
	}
	return ERC20Json, TokenVestingJson, nil
}

func TestDeployERC20(t *testing.T) {
	amount := new(big.Int).SetUint64(uint64(100000))
	decimal := uint8(8)
	contentERC20, err := ioutil.ReadFile(test.Erc20JsonLocation)
	if err != nil {
		t.Error(err)
	}
	var erc20Contract ContractAbi
	err = json.Unmarshal(contentERC20, &erc20Contract)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	contentTokenVesting, err := ioutil.ReadFile(test.TokenVestingJsonLocation)
	if err != nil {
		t.Error(err)
	}
	var tokenVestingContract ContractAbi
	err = json.Unmarshal(contentTokenVesting, &tokenVestingContract)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	contractERC20, err := test.Web3.Thk.NewContract(erc20Contract.Abi)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	contractTokenVesting, err := test.Web3.Thk.NewContract(tokenVestingContract.Abi)
	from := test.TmpAddress
	nonce, err := test.Web3.Thk.GetNonce(from, test.Web3.Thk.DefaultChainId)
	if err != nil {
		t.Error("get nonce error", err)
		t.FailNow()
	}
	transaction := util.Transaction{
		ChainId: test.Web3.Thk.DefaultChainId, FromChainId: test.Web3.Thk.DefaultChainId, ToChainId: test.Web3.Thk.DefaultChainId, From: from,
		To: "", Value: "0", Input: "", Nonce: strconv.Itoa(int(nonce)),
	}

	hashERC20, err := contractERC20.Deploy(transaction, erc20Contract.Bytecode, test.TmpKey, Symbol, Name, decimal, amount)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("contractERC20 hash:", hashERC20)

	time.Sleep(time.Second * 10)
	receiptERC20, err := test.Web3.Thk.GetTransactionByHash(test.Web3.Thk.DefaultChainId, hashERC20)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("contract addr:", receiptERC20.ContractAddress)
	toERC20 := receiptERC20.ContractAddress
	newtransaction := util.Transaction{
		ChainId: test.Web3.Thk.DefaultChainId, FromChainId: test.Web3.Thk.DefaultChainId, ToChainId: test.Web3.Thk.DefaultChainId, From: from,
		To: toERC20, Value: "0", Input: "", Nonce: strconv.Itoa(int(nonce)),
	}
	result, err := contractERC20.Call(newtransaction, "symbol", test.TmpKey)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("result:", result)
	var str string
	err = contractERC20.Parse(result.Out, "symbol", &str)
	if err != nil {
		t.Error(err)
		println("failed")
		return
	}
	//TokenVesting
	nonce, err = test.Web3.Thk.GetNonce(from, test.Web3.Thk.DefaultChainId)
	if err != nil {
		t.Error(err)
		println("get nonce error")
		return
	}
	deployTokenVesting := util.Transaction{
		ChainId: test.Web3.Thk.DefaultChainId, FromChainId: test.Web3.Thk.DefaultChainId, ToChainId: test.Web3.Thk.DefaultChainId, From: from,
		To: "", Value: "0", Input: "", Nonce: strconv.Itoa(int(nonce)),
	}
	toERC20Bytes, err := hexutil.Decode(toERC20)
	addressERC20 := common.BytesToAddress(toERC20Bytes)
	hashTokenVesting, err := contractTokenVesting.Deploy(deployTokenVesting, tokenVestingContract.Bytecode, test.TmpKey, addressERC20)
	if err != nil {
		t.Error(err)
		fmt.Println("get contractTokenVesting hash error")
		return
	}
	t.Log("contractTokenVesting hash:", hashTokenVesting)
	time.Sleep(time.Second * 10)
	receiptTokenVesting, err := test.Web3.Thk.GetTransactionByHash(test.Web3.Thk.DefaultChainId, hashTokenVesting)
	if err != nil {
		t.Error(err)
		fmt.Println("get hash error")
		return
	}
	t.Log("contractTokenVesting addr:", receiptTokenVesting.ContractAddress)
	toTokenVesting := receiptTokenVesting.ContractAddress

	//Approve
	addrefrom, err := hexutil.Decode(from)
	addressfrom := common.BytesToAddress(addrefrom)
	value := new(big.Int).SetUint64(uint64(100000))
	input, err := contractERC20.GetInput("approve", addressfrom, value)
	println(input)

	nonce, err = test.Web3.Thk.GetNonce(from, test.Web3.Thk.DefaultChainId)
	if err != nil {
		t.Error(err)
		println("get nonce error")
		return
	}
	transactionApprove := util.Transaction{
		ChainId: test.Web3.Thk.DefaultChainId, FromChainId: test.Web3.Thk.DefaultChainId, ToChainId: test.Web3.Thk.DefaultChainId, From: from,
		To: toERC20, Value: "0", Input: input,
		Nonce: strconv.Itoa(int(nonce)),
	}
	err = test.Web3.Thk.SignTransaction(&transactionApprove, test.TmpKey)

	txhash, err := test.Web3.Thk.SendTx(&transactionApprove)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	println(txhash)
	t.Log("Approve txhash:", txhash)

	time.Sleep(time.Second * 10)

	//Transfer
	addr := strings.ToLower(toTokenVesting)
	addreto, err := hexutil.Decode(addr)
	addressto := common.BytesToAddress(addreto)
	value = new(big.Int).SetUint64(uint64(100000))
	input, err = contractERC20.GetInput("transferFrom", addressfrom, addressto, value)
	//input, err := contract.GetInput("transfer",  addressto, value)
	println(input)
	nonce, err = test.Web3.Thk.GetNonce(from, test.Web3.Thk.DefaultChainId)
	if err != nil {
		t.Error(err)
		println("get nonce error")
		return
	}
	transactionTransfer := util.Transaction{
		ChainId: test.Web3.Thk.DefaultChainId, FromChainId: test.Web3.Thk.DefaultChainId, ToChainId: test.Web3.Thk.DefaultChainId, From: from,
		To: toERC20, Value: "0", Input: input,
		Nonce: strconv.Itoa(int(nonce)),
	}
	err = test.Web3.Thk.SignTransaction(&transactionTransfer, test.TmpKey)

	txhash, err = test.Web3.Thk.SendTx(&transactionTransfer)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("Transfer txhash:", txhash)

	time.Sleep(time.Second * 10)
	//Vesting
	nonce, err = test.Web3.Thk.GetNonce(from, test.Web3.Thk.DefaultChainId)
	if err != nil {
		t.Error(err)
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
	tmcliffP := new(big.Int).SetInt64(tmcliff.Unix())
	tmstartP := new(big.Int).SetInt64(tmstart.Unix())
	tmendP := new(big.Int).SetInt64(tmend.Unix())
	timesP := new(big.Int).SetInt64(2)
	total := new(big.Int).SetUint64(uint64(100000))
	resultTokenVesting, err := contractTokenVesting.Send(transactionTokenVesting, "addPlan", test.TmpKey,
		vestAddress, tmcliffP, tmstartP, timesP, tmendP, total, false, "上交易所后私募锁仓10分钟，之后每10分钟释放50%")
	if err != nil {
		t.Error(err)
	}
	t.Log("result:", resultTokenVesting)

	var token Token
	token.Name = Name
	token.Symbol = Symbol
	token.Total = 100000
	token.ContractAddress = receiptERC20.ContractAddress
	token.ABI = erc20Contract.Abi
	token.Icon = "icon"
	token.Website = "www.thinkey.com"
	token.Introduction = Des
	token.Date = time.Now().Format("2006-01-02")
	token.ChainId = "2"
	token.Decimal = 8
	PostInfo(token)
	//PostFile(token, "../resources/abc.sol")
}
