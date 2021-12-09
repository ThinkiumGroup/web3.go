package common

import (
	"encoding/json"
	"fmt"
	"github.com/ThinkiumGroup/web3.go/common/hexutil"
	"github.com/ThinkiumGroup/web3.go/web3/thk/util"
	"testing"
)

type Tx struct {
	Tx   string
	Hash string
}

func TestHash(t *testing.T) {
	var txs []Tx
	txs = append(txs, Tx{
		Tx: `{
        "chainid": "103",
        "from": "0x7112add80e015c16c84f2e7cdb52c4f70bcb2e60",
        "to": "0x33819ba73fb9a63b547815822c044530124bd4b1",
        "nonce": "74",
        "value": "0",
        "input": "0xa9059cbb0000000000000000000000003d5c0a7f100d2b17acd88accd533c3b1512fdd520000000000000000000000000000000000000000000000000000000000000001",
        "hash": "",
        "uselocal": false,
        "extra": "0x",
        "timestamp": 0
    }`,
		Hash: "0x20773446611df2c399cd2aff52f50cb5d0973a560695e44bd8df5873d6d998a4",
	})
	txs = append(txs, Tx{
		Tx: `{
        "chainid": "1",
        "from": "0x07f9aaabd576b732486208e21577e01557c297d7",
        "to": "0x8b95043faa2aa106505a6e133c034ef6003d23a4",
        "nonce": "119",
        "value": "17067600000000000000",
        "input": "0x",
        "hash": "",
        "uselocal": false,
        "extra": "0x",
        "timestamp": 0
    }`,
		Hash: "0x8abcc013b06dd2ba5bfa53f26143e04d7365a877bc236fbc0b73d3121aa9a1c6",
	})
	for _, tx := range txs {
		var transaction util.Transaction
		err := json.Unmarshal([]byte(tx.Tx), &transaction)
		if err != nil {
			t.Error(err.Error())
			return
		}
		hashBytes, err := transaction.HashValue()
		if err != nil {
			t.Error(err.Error())
			return
		}
		realHash := hexutil.Encode(hashBytes)
		fmt.Println(realHash)
		fmt.Println(realHash == tx.Hash)
	}
}
