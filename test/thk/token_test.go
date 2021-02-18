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
	"web3.go/web3/thk"
	"web3.go/web3/thk/util"
)

type ContractJson struct {
	ContractName string `json:"contractName"`
	ABI          string `json:"abi"`
	ByteCode     string `json:"bytecode"`
}

const (
	from                     = test.TmpAddress
	chainId                  = "1"
	Symbol                   = "USDT"
	Name                     = "Token of USD"
	decimals                 = uint8(8)
	Erc20JsonLocation        = "../resources/ERC20.json"
	TokenVestingJsonLocation = "../resources/TokenVesting.json"
	erc20Address             = "0x211df4ab43bab418f1fa82a15dcca5ce9dab15a9"
	tokenVestingAddress      = "0x9490e52e718bba165b3323a53316fc7be5942964"
)

var (
	erc20Bytecode, tokenVestingBytecode string
	tokenVesting, erc20                 *thk.Contract
	fromAddress                         common.Address
	totalSupply                         = new(big.Int).SetUint64(uint64(1000000000))
	approveAmount                       = new(big.Int).SetUint64(uint64(100000000))
)

type ContractAbi struct {
	Abi      string `json:"abi"`
	Bytecode string `json:"bytecode"`
}

func init() {
	var err error
	erc20, erc20Bytecode, err = loadContract(Erc20JsonLocation)
	if err != nil {
		panic(err)
	}
	tokenVesting, tokenVestingBytecode, err = loadContract(TokenVestingJsonLocation)
	if err != nil {
		panic(err)
	}

	toAddress, err := hexutil.Decode(test.TmpAddress)
	if err != nil {
		panic(err)
	}
	fromAddress = common.BytesToAddress(toAddress)
}
func TestCompileERC20(t *testing.T) {
	contractMap, err := compiler.CompileContract("../resources/contract/ERC20.sol",
		"../resources/contract/IERC20.sol", "../resources/contract/Pausable.sol",
		"../resources/contract/SafeMath.sol", "../resources/contract/Ownable.sol",
		"../resources/contract/Address.sol", "../resources/contract/SafeERC20.sol",
		"../resources/contract/TokenVesting.sol")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	var ERC20Json, TokenVestingJson ContractJson
	for k, contract := range contractMap {
		arr := strings.Split(k, ":")
		if arr[len(arr)-1] == "ERC20" {
			ERC20Json.ByteCode = contract.Code
			abi := contract.Info.AbiDefinition
			abiBytes, _ := json.Marshal(abi)
			ERC20Json.ABI = string(abiBytes)
		}
		if arr[len(arr)-1] == "TokenVesting" {
			ERC20Json.ByteCode = contract.Code
			abi := contract.Info.AbiDefinition
			abiBytes, _ := json.Marshal(abi)
			TokenVestingJson.ABI = string(abiBytes)
		}
	}
	bytes, err := json.MarshalIndent(ERC20Json, "", "  ")
	if ioutil.WriteFile(Erc20JsonLocation, bytes, 0644) == nil {
		fmt.Println("write success")
	}
	bytes, err = json.MarshalIndent(TokenVestingJson, "", "  ")
	if ioutil.WriteFile(TokenVestingJsonLocation, bytes, 0644) == nil {
		fmt.Println("write success")
	}
}

func TestERC20Deploy(t *testing.T) {
	nonce, err := test.Web3.Thk.GetNonce(from, chainId)
	if err != nil {
		t.Error("get nonce error", err)
		t.FailNow()
	}
	transaction := util.Transaction{
		ChainId: chainId, FromChainId: chainId, ToChainId: chainId, From: from,
		To: "", Value: "0", Input: "", Nonce: strconv.Itoa(int(nonce)),
	}

	hash, err := erc20.Deploy(transaction, erc20Bytecode, test.TmpKey, Symbol, Name, decimals, totalSupply)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("erc20 hash:", hash)

	receipt, err := test.BlockGetTransactionReceipt(chainId, hash)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if receipt.Status == 0 {
		t.Error(fmt.Sprintf("failed tx, hash:%v", hash))
		t.FailNow()
	}
	erc20Address := receipt.ContractAddress
	t.Log("erc20Address addr:", erc20Address)
	var symbol string
	err = erc20.CallAndParse(chainId, erc20Address, &symbol, "symbol")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("result:", symbol)
}

// deploy token-vesting
func TestTokenVestingDeploy(t *testing.T) {
	bytes, err := hexutil.Decode(erc20Address)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	addressERC20 := common.BytesToAddress(bytes)
	nonce, err := test.Web3.Thk.GetNonce(from, chainId)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	transaction := util.Transaction{
		ChainId: chainId, FromChainId: chainId, ToChainId: chainId, From: from,
		To: "", Value: "0", Input: "", Nonce: strconv.Itoa(int(nonce)),
	}
	hash, err := tokenVesting.Deploy(transaction, tokenVestingBytecode, test.TmpKey, addressERC20)
	if err != nil {
		t.Error(err)
		fmt.Println("get tokenVesting hash error")
		return
	}
	t.Log("tokenVesting hash:", hash)
	receipt, err := test.BlockGetTransactionReceipt(chainId, hash)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if receipt.Status == 0 {
		t.Error(fmt.Sprintf("failed tx, hash:%v", hash))
		t.FailNow()
	}
	tokenVestingAddress := receipt.ContractAddress
	t.Log("tokenVestingAddress:", tokenVestingAddress)
}

//////////////////////////////////////////////////////////////////////////////////////////

// Approve
func TestErc20Approve(t *testing.T) {
	input, err := erc20.GetInput("approve", fromAddress, approveAmount)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("Approve input:", input)

	nonce, err := test.Web3.Thk.GetNonce(from, chainId)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	transaction := util.Transaction{
		ChainId: chainId, FromChainId: chainId, ToChainId: chainId, From: from,
		To: erc20Address, Value: "0", Input: input, Nonce: strconv.Itoa(int(nonce)),
	}
	err = test.Web3.Thk.SignTransaction(&transaction, test.TmpKey)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	hash, err := test.Web3.Thk.SendTx(&transaction)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("Approve hash:", hash)
	checkTransaction(t, hash)
}

//Transfer
func TestErc20Transfer(t *testing.T) {
	bytes, err := hexutil.Decode(strings.ToLower(tokenVestingAddress))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	addressTokenVesting := common.BytesToAddress(bytes)
	input, err := erc20.GetInput("transferFrom", fromAddress, addressTokenVesting, approveAmount)
	//input, err := erc20.GetInput("transfer", addressTokenVesting, value)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("transferFrom input:", input)
	nonce, err := test.Web3.Thk.GetNonce(from, chainId)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	transaction := util.Transaction{
		ChainId: chainId, FromChainId: chainId, ToChainId: chainId, From: from,
		To: erc20Address, Value: "0", Input: input, Nonce: strconv.Itoa(int(nonce)),
	}
	err = test.Web3.Thk.SignTransaction(&transaction, test.TmpKey)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	hash, err := test.Web3.Thk.SendTx(&transaction)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("Transfer hash:", hash)

	checkTransaction(t, hash)
}

// token-vest addPlan
func TestTokenVestingAddPlan(t *testing.T) {
	TestErc20Approve(t)
	TestErc20Transfer(t)
	startTime := new(big.Int).SetInt64(time.Now().Unix())
	t.Log("startTime:", startTime)
	lockToTime := new(big.Int).SetInt64(time.Now().Unix() + 2*60)
	t.Log("lockToTime:", lockToTime)
	endTime := new(big.Int).SetInt64(time.Now().Unix() + 5*60)
	t.Log("endTime:", endTime)
	releaseStages := new(big.Int).SetInt64(2)
	total := new(big.Int).SetUint64(uint64(100000))
	remark := "After going to the stock exchange, private placement lock positions for 10 minutes, and then release 50% every 10 minutes"
	nonce, err := test.Web3.Thk.GetNonce(from, chainId)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	transaction := util.Transaction{
		ChainId: chainId, FromChainId: chainId, ToChainId: chainId, From: from,
		To: tokenVestingAddress, Value: "0", Input: "", Nonce: strconv.Itoa(int(nonce)),
	}
	hash, err := tokenVesting.Send(transaction, "addPlan", test.TmpKey, fromAddress, startTime, lockToTime, releaseStages, endTime, total, false, remark)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("addPlan hash:", hash)
	checkTransaction(t, hash)
}

// release
func TestTokenVestingRelease(t *testing.T) {
	var releasableAmount *big.Int
	err := tokenVesting.CallAndParse(chainId, tokenVestingAddress, &releasableAmount, "releasableAmount", fromAddress)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("releasableAmount:", releasableAmount)
	var vestedAmount *big.Int
	err = tokenVesting.CallAndParse(chainId, tokenVestingAddress, &vestedAmount, "vestedAmount", fromAddress)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("vestedAmount:", vestedAmount)

	nonce, err := test.Web3.Thk.GetNonce(from, chainId)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	transaction := util.Transaction{
		ChainId: chainId, FromChainId: chainId, ToChainId: chainId, From: from,
		To: tokenVestingAddress, Value: "0", Input: "", Nonce: strconv.Itoa(int(nonce)),
	}
	hash, err := tokenVesting.Send(transaction, "release", test.TmpKey, fromAddress)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("release.hash:", hash)
	checkTransaction(t, hash)
}

func checkTransaction(t *testing.T, hash string) {
	err := test.BlockCheckTransactionReceipt(chainId, hash)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func loadContract(jsonFileLocation string) (*thk.Contract, string, error) {
	file, err := ioutil.ReadFile(jsonFileLocation)
	if err != nil {
		return nil, "", err
	}
	var contractAbi ContractAbi
	err = json.Unmarshal(file, &contractAbi)
	if err != nil {
		return nil, "", err
	}
	contract, err := test.Web3.Thk.NewContract(contractAbi.Abi)
	if err != nil {
		return nil, "", err
	}
	return contract, contractAbi.Bytecode, nil
}
