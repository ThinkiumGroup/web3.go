package test

import (
	"fmt"
	"strconv"
	"testing"
	"web3.go/test"
	"web3.go/web3/thk/util"
)

func TestSendTx(t *testing.T) {
	var err error
	to := test.TmpAddress
	nonce, err := test.Web3.Thk.GetNonce(test.Web3.Thk.DefaultAddress, "1")
	fmt.Println("nonce:", nonce)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	transaction := util.Transaction{
		ChainId: test.Web3.Thk.DefaultChainId, FromChainId: test.Web3.Thk.DefaultChainId, ToChainId: test.Web3.Thk.DefaultChainId, From: test.Web3.Thk.DefaultAddress,
		To: to, Value: test.DefaultValue, Input: "", Nonce: strconv.Itoa(int(nonce)), UseLocal: false, Extra: "",
	}
	err = test.Web3.Thk.SignTransaction(&transaction, test.Web3.Thk.DefaultPrivateKey)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	hash, err := test.Web3.Thk.SendTx(&transaction)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("hash:", hash)
}

func TestGetTransactionByHash(t *testing.T) {
	var err error
	hash := "0x206ba760f935e6d2f7f2ad8ee776ed07b5e4a9ea6948a40b6fa48b08ea75b957"
	res, err := test.Web3.Thk.GetTransactionByHash("1", hash)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("res:%+v", test.JsonFormat(res))
}
