package test

import (
	"fmt"
	"strconv"
	"testing"
	"time"
	"web3.go/test"
	"web3.go/web3/thk"
	"web3.go/web3/thk/util"
)

var (
	fromChainId = "1"
	toChainId   = "2"
	chequeValue = "6" + "000000000000000000"
	expireAfter = 200
)

//Cross chain transfer process, from a chain to B chain
//1 - first, write a check from chain a, send a transaction to [system cross chain withdrawal contract (0x000000000000000)], and specify the overdue height of the check -- that is, when the height of chain B exceeds (excluding) this value, the check cannot be withdrawn, but can only be returned
//2 - check the previous transaction result
//3 - if the transaction is successful, obtain the check deposit certificate of chain B from chain a
//4 - send a transaction from chain B to [system cross chain deposit contract (0x000000000003000000)] with the returned check certificate as input
//5 - check the transaction result of the previous step. If the transaction is successful, the cross chain transfer is successful
//Cancel check - the check cannot be cancelled until the specified block height is reached
//6 - if the check reaches the specified block height and is not withdrawn, the cross chain deposit needs to be cancelled manually
//7 - get the check cancellation certificate from the B chain and use it as input to send a transaction to [system cross chain cancellation deposit contract (0x000000000003000000)]
//8 - check the transaction status of the previous step. If it fails, try step 6, 7, 8 again or contact inter core technology for assistance

func TestEncodeDecode(t *testing.T) {
	input := "0x00000001f167a1c5c5fab6bddca66118216817af3fa86827000000000000018500000067f167a1c5c5fab6bddca66118216817af3fa8682700000000005ef43c20000000000000000000000000000000000000000000000001158e460913d00000"
	var cash thk.CashCheque
	err := cash.Decode(input)
	if err != nil {
		t.Log(err)
	}
	fmt.Println(test.JsonFormat(cash))
}

func TestTransferAcrossChain(t *testing.T) {
	expireAfter = 200
	fmt.Println("===Write a check===")
	cheque := genCheque(t)
	fmt.Println("===Get proof of cashing a check===")
	proof := getChequeProof(cheque, t)

	fmt.Println("===Cash a check===")
	tx := util.Transaction{
		ChainId: toChainId, FromChainId: toChainId, ToChainId: toChainId, From: test.Web3.Thk.DefaultAddress,
		To: thk.SystemContractAddressDeposit, Value: "0", Input: proof,
	}
	cashOrCancelCheque(tx, t)
}

func TestCancelCheque(t *testing.T) {
	expireAfter = 3
	fmt.Println("===Write a check===")
	cheque := genCheque(t)
	expireHeight, _ := strconv.Atoi(cheque.ExpireHeight)
	fmt.Println("===Waiting for the check to expire===")
	for {
		time.Sleep(2 * time.Second)
		stats, err := test.Web3.Thk.GetStats(toChainId)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		fmt.Printf("ExpireHeight: %v, CurrentHeight: %v\n", cheque.ExpireHeight, stats.CurrentHeight)
		if stats.CurrentHeight > expireHeight {
			break
		}
	}

	fmt.Println("===The check has expired. Generate proof to cancel the check===")
	time.Sleep(2 * time.Second)
	proofCancel := getCancelChequeProof(cheque, t)
	tx2 := util.Transaction{
		ChainId: fromChainId, FromChainId: fromChainId, ToChainId: fromChainId, From: test.Web3.Thk.DefaultAddress,
		To: thk.SystemContractAddressCancel, Value: "0", Input: proofCancel,
	}
	fmt.Println(test.JsonFormat(tx2))
	fmt.Println("===Cancel a check===")
	cashOrCancelCheque(tx2, t)
}

func genCheque(t *testing.T) *thk.CashCheque {
	chainInfo, err := test.Web3.Thk.GetStats(toChainId)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	chainInfo.ChainId = 2
	fmt.Printf("chainInfo:%+v\n", chainInfo)
	expireHeight := strconv.Itoa(chainInfo.CurrentHeight + expireAfter)

	nonce, err := test.Web3.Thk.GetNonce(test.Web3.Thk.DefaultAddress, fromChainId)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	cashCheque := &thk.CashCheque{
		ChainId:      fromChainId,
		FromChainId:  fromChainId,
		From:         test.Web3.Thk.DefaultAddress,
		Nonce:        strconv.Itoa(int(nonce)),
		ToChainId:    toChainId,
		To:           test.Web3.Thk.DefaultAddress,
		ExpireHeight: expireHeight,
		Value:        chequeValue,
	}
	fmt.Printf("cashCheque.Nonce: %v\n", cashCheque.Nonce)
	chequeAsInput, err := cashCheque.Encode()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Printf("chequeAsInput: %v\n", chequeAsInput)
	transaction := util.Transaction{
		ChainId: fromChainId, FromChainId: fromChainId, ToChainId: toChainId, From: test.Web3.Thk.DefaultAddress,
		To: thk.SystemContractAddressWithdraw, Value: "0", Input: chequeAsInput, Nonce: strconv.Itoa(int(nonce)),
	}
	err = test.Web3.Thk.SignTransaction(&transaction, test.Web3.Thk.DefaultPrivateKey)
	hash, err := test.Web3.Thk.SendTx(&transaction)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("hash:%v\n", hash)

	time.Sleep(5 * time.Second)
	res, err := test.Web3.Thk.GetTransactionByHash(fromChainId, hash)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("txRes:%+v\n", res)
	if res.Status != 1 {
		t.Error("gen cheque error")
		t.FailNow()
	}
	return cashCheque
}

func getChequeProof(cashCheque *thk.CashCheque, t *testing.T) string {
	time.Sleep(5 * time.Second)
	proofRes, err := test.Web3.Thk.RpcMakeVccProof(cashCheque)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("proofRes:%+v\n", proofRes)
	if proofRes["errCode"] != nil {
		t.Error("get vccProofInput error")
		t.FailNow()
	}
	return proofRes["input"].(string)
}

func getCancelChequeProof(cashCheque *thk.CashCheque, t *testing.T) string {
	cashCheque.ChainId = toChainId
	res, err := test.Web3.Thk.MakeCCCExistenceProof(cashCheque)
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	if res["errCode"] != nil {
		t.Error(res["errMsg"])
		t.FailNow()
	}
	fmt.Printf("res:%v\n", res)
	if res["existence"].(bool) == true {
		t.Error("cheque has been cashed")
		t.FailNow()
	}
	return res["input"].(string)
}

func cashOrCancelCheque(tx util.Transaction, t *testing.T) {
	nonce, err := test.Web3.Thk.GetNonce(test.Web3.Thk.DefaultAddress, tx.ChainId)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	tx.Nonce = strconv.Itoa(int(nonce))

	err = test.Web3.Thk.SignTransaction(&tx, test.Web3.Thk.DefaultPrivateKey)
	hash, err := test.Web3.Thk.SendTx(&tx)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("hash:%v\n", hash)

	time.Sleep(5 * time.Second)
	res, err := test.Web3.Thk.GetTransactionByHash(tx.ChainId, hash)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("txRes:%+v\n", res)
	if res.Status != 1 {
		t.Error("tx error")
		t.FailNow()
	}
}
