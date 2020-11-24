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

// 跨链转账流程，从A链跨链转到B链
// 1-首先从A链开支票，发送一笔交易到【系统跨链取款合约(0x0000000000000000000000000000000000020000)】, 需要指定支票过期高度--指的是当B链高度超过（不含）这个值时，这张支票不能被支取，只能退回
// 2-检查上一步交易结果
// 3-如果交易成功，从A链获取B链支票存款证明
// 4-将返回的支票证明作为input，从B链发送一笔交易到【系统跨链存款合约(0x0000000000000000000000000000000000030000)】
// 5-检查上步骤交易结果，如果交易成功则跨链转账成功

// 撤销支票-达到指定块高之后才能被取消
// 6-如果到达指定块高支票未被支取则需要手动取消跨链存款
// 7-从B链获取支票撤销证明并将之作为input，发送一笔交易到【系统跨链撤销存款合约(0x0000000000000000000000000000000000030000)】
// 8-检查上步骤交易状态，如果失败可重试6、7、8步骤或者联系芯际技术协助处理

func TestEncodeDecode(t *testing.T) {
	input := "0x00000001f167a1c5c5fab6bddca66118216817af3fa86827000000000000018500000067f167a1c5c5fab6bddca66118216817af3fa8682700000000005ef43c20000000000000000000000000000000000000000000000001158e460913d00000"
	var cash thk.CashCheque
	err := cash.Decode(input)
	if err != nil {
		t.Log(err)
	}
	fmt.Println(test.JsonFormat(cash))
}

//跨链转账测试
func TestTransferAcrossChain(t *testing.T) {
	expireAfter = 200
	fmt.Println("===开支票===")
	cheque := genCheque(t)
	fmt.Println("===获取兑现支票的证明===")
	proof := getChequeProof(cheque, t)

	fmt.Println("===兑现支票===")
	tx := util.Transaction{
		ChainId: toChainId, FromChainId: toChainId, ToChainId: toChainId, From: test.Web3.Thk.DefaultAddress,
		To: thk.SystemContractAddressDeposit, Value: "0", Input: proof,
	}
	cashOrCancelCheque(tx, t)
}

//跨链转账过期取消支票
func TestCancelCheque(t *testing.T) {
	expireAfter = 3
	fmt.Println("===开支票===")
	cheque := genCheque(t)
	expireHeight, _ := strconv.Atoi(cheque.ExpireHeight)
	fmt.Println("===等待支票过期===")
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

	fmt.Println("===支票已过期，生成撤销支票的证明===")
	time.Sleep(2 * time.Second)
	proofCancel := getCancelChequeProof(cheque, t)
	tx2 := util.Transaction{
		ChainId: fromChainId, FromChainId: fromChainId, ToChainId: fromChainId, From: test.Web3.Thk.DefaultAddress,
		To: thk.SystemContractAddressCancel, Value: "0", Input: proofCancel,
	}
	fmt.Println(test.JsonFormat(tx2))
	fmt.Println("===撤销支票===")
	cashOrCancelCheque(tx2, t)
}

//开支票
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

//获取兑现支票的证明
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

//生成撤销支票的证明
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

//兑现或取消支票
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
