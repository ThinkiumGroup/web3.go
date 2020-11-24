package common

import (
	"fmt"
	"testing"
	"web3.go/common"
	"web3.go/common/hexutil"
)

func TestKey(t *testing.T) {
	key, _ := common.Cipher.GenerateKey()
	fmt.Println("key:", common.Bytes2Hex(key.ToBytes()))
	publicKey := key.GetPublicKey()
	fmt.Println("pub:", common.Bytes2Hex(publicKey.ToBytes()))
	fmt.Println("addr:", common.Bytes2Hex(publicKey.ToAddress()))
}

func TestSignAndVerify(t *testing.T) {
	fmt.Println("======sign and verify======")
	msg := "123"
	sig, _ := common.Sign(msg, "7b3effbc3292e156d1993f8327e6e5d9fe776a5494bc911baee53aa1db0be6d6")
	sigBytes, _ := hexutil.Decode(sig)
	pub, _ := hexutil.Decode("0x04fa8cc34e0ba0f701fa256c58255846cfa4bf529bf144ac226d8f2df8e9019d7ed48b979df426fe82609b77f0c582e98aeb587b4392a8f9f7447fa294a59c7548")
	msgHash := common.Hash256(msg)
	verify := common.Cipher.Verify(pub, msgHash, sigBytes)
	fmt.Println("verify:", verify)
}
