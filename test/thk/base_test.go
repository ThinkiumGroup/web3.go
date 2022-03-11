package test

import (
	"fmt"
	"github.com/ThinkiumGroup/go-common"
	"github.com/ThinkiumGroup/web3.go/common/hexutil"
	"github.com/ThinkiumGroup/web3.go/test"
	"strconv"
	"testing"
)

func TestThkGetStats(t *testing.T) {
	stats, err := test.Web3.Thk.GetStats("1")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	fmt.Printf("stats:%+v", stats)
}
func TestThkGetChainInfo(t *testing.T) {
	var chainIds = []int{}
	infos, err := test.Web3.Thk.GetChainInfo(chainIds)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("infos:%+v", infos)
}

func TestThkGetBlockHeader(t *testing.T) {
	res, err := test.Web3.Thk.GetBlockHeader("1", "983918")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("res:%+v", res)
}

func TestThkGetBlock(t *testing.T) {
	res, err := test.Web3.Thk.GetBlock("1", "983918")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for _, pas := range res.BlockPass {
		if pas == nil {
			continue
		}
		t.Logf("pub:%s\nsig: %s", hexutil.Encode(pas.PublicKey.Bytes()), hexutil.Encode(pas.Signature.Bytes()))
		h := res.BlockHeader.Hash()
		t.Log(h.String())
		if !common.VerifyMsg(res.BlockHeader, pas.PublicKey.Bytes(), pas.Signature.Bytes()) {
			t.Log("check faild")
			t.FailNow()
		}
	}

	fmt.Printf("res:%+v", res)
}

func TestGetTxProof(t *testing.T) {
	res, err := test.Web3.Thk.GetTxProof("1", "0x657bf5ab9e1f51100d31c5b049b8abac159359135ed04bd6b47ce9af116f0221")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("res:%+v", test.JsonFormat(res))
}

func TestThkPing(t *testing.T) {
	res, err := test.Web3.Thk.Ping("192.168.1.14:23024")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("res:%+v", res)
}

func TestThkGetBlockTxs(t *testing.T) {
	res, err := test.Web3.Thk.GetBlockTxs("0", "3613", "1", "10")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("res:%+v", res)
}

func TestThkGetAccount(t *testing.T) {
	account, err := test.Web3.Thk.GetAccount(test.Web3.Thk.DefaultAddress, "1")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("account:", account)
}
func TestThkGetBalance(t *testing.T) {
	bal, err := test.Web3.Thk.GetBalance(test.Web3.Thk.DefaultAddress, "1")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("Balance:", bal)
}

func TestThkGetNonce(t *testing.T) {
	nonce, err := test.Web3.Thk.GetNonce(test.Web3.Thk.DefaultAddress, "1")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("nonce:", nonce)
}

func TestThkGetCommittee(t *testing.T) {
	var err error
	res, err := test.Web3.Thk.GetCommittee("0", "1")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("res", res)
}

func TestThkGetNodeSig(t *testing.T) {
	var err error
	test.Web3.Thk.DefaultAddress = "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23"
	nonce, err := test.Web3.Thk.GetNonce(test.Web3.Thk.DefaultAddress, "1")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	privateKey := "0x15d115381a4e445d66c59f4c2b884d78a34ac54bccc333b4508bce9cacf32539"
	nodeId := "0x3d85aa2a649fa5fd421988cebff58d7173f7b563b8a9594e92bcf3e9f5e43037c3463121af51aacc8a8cf2d8cfcc6fa717b774fc0aceec04d7185c87e279c1f6"
	res, err := test.Web3.Thk.GetNodeSig(nodeId, "1", "0x2c7536e3605d9c16a7a3d7b1898e529396a65c23", strconv.FormatInt(nonce, 10), "5000", privateKey)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("NodeSig:", res)
}
