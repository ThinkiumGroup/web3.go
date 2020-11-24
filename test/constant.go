package test

import (
	"encoding/json"
	"web3.go/web3"
	"web3.go/web3/providers"
)

var (
	RpcHost                  = "rpctest.thinkium.org"
	Web3                     = web3.NewWeb3(providers.NewHTTPProvider(RpcHost, 10, false))
	DefaultValue             = "1" + "000000000000000000"
	TmpKey                   = "0xc614545a9f1d9a2eeda26836e42a4c11631f25dc3d0dcc37fe62a89c4ff293d1"
	TmpAddress               = "0x5dfcfc6f4b48f93213dad643a50228ff873c15b9"
	Erc20JsonLocation        = "../resources/ERC20.json"
	TokenVestingJsonLocation = "../resources/TokenVesting.json"
)

func init() {
	Web3.Thk.DefaultPrivateKey = "0x8e5b44b6cee8fa05092b4b5a8843aa6b0ec37915a940c9b5938e88a7e6fdd83a"
	Web3.Thk.DefaultAddress = "0xf167a1c5c5fab6bddca66118216817af3fa86827"
	Web3.Thk.DefaultChainId = "1"
}
func JsonFormat(obj interface{}) string {
	jsonStr, _ := json.MarshalIndent(obj, "", "\t")
	return string(jsonStr)
}
