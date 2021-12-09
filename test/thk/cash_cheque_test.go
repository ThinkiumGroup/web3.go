package test

import (
	"fmt"
	"github.com/ThinkiumGroup/web3.go/test"
	"github.com/ThinkiumGroup/web3.go/web3/thk"
	"github.com/ThinkiumGroup/web3.go/web3/thk/util"
	"strconv"
	"testing"
	"time"
)

var (
	fromChainId = "1"
	toChainId   = "2"
	chequeValue = "1" + "000000000000000000"
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
	input := "0x00000001700fe44d941225d58e695c449f79412cc7fdbcf8000000000000000500000002700fe44d941225d58e695c449f79412cc7fdbcf80000000000cd35942000000000000000000000000000000000000000000000080900a88ac615440000"
	cashInput := "0x9500000001700fe44d941225d58e695c449f79412cc7fdbcf8000000000000000800000002700fe44d941225d58e695c449f79412cc7fdbcf80000000000cd4ab00a0808bb44f943d050000000a40115c15bc00a6ecc6c41df074d3b206b4c92bf0135eaa1ad2e1bc5cb8ea19b5c0d90b5faf297941093a1b1de98b98e55dcd716559a4637f8d1093d3d2094acfb23a22d4a6beed4b8fc01c200008080940b934080c2ad4b80810004873886e8ed6d1fb4cfb76ee3a51743017fb843d4d75aebd1c7ff74af54a11efe6fe278546f23095d9045730bf9754685dc95a408c4c7187b75dd5c56ae2e4e40287dd0ee81fb16613e9c19f51b9ca42dd34bc08339a27817d402f7959b13f5235308633a3ff8f0cb4de954288da3be5873100efe659c25ef05cd6bc0653e7257000107940e934080c2ffff80810004511655f6352ed0decc7276b684df10bb46feff9a908acf0df2bd6016a17ec973a03f7a59bd2949c4f5d4a38522e7ca2a3c2f4734cbf3303a6cf40f5d037a25016ec88920b0aed933b26975be441f071e7bd8f93e78d5ca1430177b67d7d5fb0c6d0e347c5679245f580a4b2da483c3019e575111b6ce65978fa59d220db2548000010e9404934080c2ffff80810004844fba978763858c37126d7c256598ea19a781b6ed81af5bd152c40d406f60fff11def743b04694881d7493ad4c3cb448f00ed0ab90e11d6a587fd99f0dc0e4986adab2d712eefb87c4c5c462a36efa0a8ea8e33103f7a0d21e7998fc582c00dce7e3cc6f7f75806eff5c6a5cc687e9c04f90079ae5b34b65a4b5cbf6f0070c40001049424930080c20000c089e7efb17aeab2d0c342d07f124096cb88aa39b6e634d560e4ae49d5cab5ee53810005423ba7e430ec3184e82c8df70aa2c5df8af08a01dae25e1331089cb7c3a29ec508119a9021bd4e60a1630ed667f901dccc3464f02ff1f203200b87d0152b1c915c3e0f832b46bcce0f17b37c08c1a268808ff789fae75d6709a702c7e1727a66d9c741a52af4683582993cf7839c8597457a0befe67682bee502d1f21e49db22db14fadeae9557e9abb544add89d2a2d22c35591f04b178f18e253f827cc03f80001109411930080c20000c012440ed40e7975972d7b55218afd52f6f7208d7cd9cd87f498717c3d77c3e5ec810002b95a5048621c6adff0e6e0501ed5f83ee6a8c195c21914dc3cff8e9ceb79b27d8f1cf6230f0498dc26f33ddb993633646822c8da316cf4574bdcb2f86872990c00009426930080c20000c0fa4cd1b148975f4b52aeb6a82f63ed2e90f9ede7c138dab58b1a135f578a959e810005b2cbeb6508244f5c9a8a0ee618068ef410d904d42bd652b6cf25576c58b04e89e2c0da358e91eed86a0658f202bb8fb2050f9e82cf7e9ba51b02bf35a7f3195f731728437fd64a11e98169b11432e9c9ebf78ca70978206e19a72ad539d82180d98df14cbc3939a40effc8c9b294d8af3733caa3370b2adc34f74c6730fd54a0590f25f3e5de0665d4004ee5514a7ab97587284f8d7ec07755a722990fba7bc9000114"
	refundInput := "0x9600000001700fe44d941225d58e695c449f79412cc7fdbcf8000000000000000500000002700fe44d941225d58e695c449f79412cc7fdbcf80000000000cd35940a080900a88ac61544000001a3cd49bdc01ed94d694f64b56225597e0cafe84a5efa709668ec8a3e4263095ff77052a7929394a1fe934080c2808280810001a09b033a1b8e6c386cf40180fafd9b9b335d9bb2ff665fa839b8de677788daac00009402934080c2ffff808100042364b4ed58f756968196789b56db2c3b1bd17b8b6044e551509d192b41db94404e7362b575d64f6c427b95366f603c9ebccbdd7bc71794fbca0c397a86836109f5bb809bb7812adf54cda134eb20eeb52224114087820bf35460afb959f3d8fc17346d5b3b64c12a0dd1bf330b2b44f9e58b9fe8350b61b70bca9f4fc96f63b2000102940a934080c2ffff80810004389db1c4fcf2cb709480a8b7ad196550c3dd7d1cc80b81057c293157f784a23fd91f1c67631704b1e5048919c2de81dc7141da9ac69bf4a4a500d71f9042e9d3d8b99f2e75c80afef28f9384917c983cc08c9c53ecdff238a309acfd6220538291017b8ffabe46fe47d45fa88aba0e238ffa40ed61cc84f06a0ee5eb285dacb300010a919425930080c20000c0c5e5be3bc5e6df278b0a01e44cdded05cfe69ee1c592e662dc6bad7ac08e4bd3810005d07655ca81f8d75de2a376801feedcf348f5ff21a1c04c45eb3b146c11456b5413d5086b86bf654eb34ab7317fe0b70109e55f6376999cad571b927be0d507222f8e331e16e453e0da085eadb2705e556b4bb4693835b3f22c1c17f1dcfccc425998e44e70d863daea2d602751182a878eaa7c6e4449325e14d92f13287df5dfea4079ecb21ff0770dc9febc35bc47160ac167362f51a0ac27edcd156d4b3da8000111"
	inputs := []string{input, cashInput, refundInput}
	for _, input := range inputs {
		var cash thk.CashCheque
		err := cash.Decode(input)
		if err != nil {
			t.Log(err)
		}
		fmt.Println(test.JsonFormat(cash))
	}
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
