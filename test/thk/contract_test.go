package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"mime/multipart"
	"net/http"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"
	"web3.go/test"
	"web3.go/test/compiler"
	"web3.go/web3/thk/util"
)

type ContractJson struct {
	ContractName string `json:"contractName"`
	ABI          string `json:"abi"`
	ByteCode     string `json:"bytecode"`
}

var Symbol = "GBDJ"
var Name = "GOLDDJ"
var Des = "金币"

type ContractAbi struct {
	Abi      string `json:"abi"`
	Bytecode string `json:"bytecode"`
}
type Token struct {
	Name            string  `json:"name"`
	Symbol          string  `json:"symbol"`
	Total           float64 `json:"total"`
	ContractAddress string  `json:"contractaddress"`
	ABI             string  `json:"abi"`
	Icon            string  `json:"icon"`
	Website         string  `json:"website"`
	Introduction    string  `json:"introduction"`
	State           string  `json:"state"`
	Date            string  `json:"date"`
	ChainId         string  `json:"chainid"`
	Decimal         int64   `json:"decimal"`
}

func TestCompile(t *testing.T) {
	ctct, err := compiler.CompileContract("../resources/contract/ERC20.sol",
		"../resources/contract/IERC20.sol", "../resources/contract/Pausable.sol",
		"../resources/contract/SafeMath.sol", "../resources/contract/Ownable.sol",
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	var ERC20Json ContractJson
	ERC20Json, _, err = GetContractJson(ctct)
	data, err := json.MarshalIndent(ERC20Json, "", "  ")
	if ioutil.WriteFile(test.Erc20JsonLocation, data, 0644) == nil {
		fmt.Println("写入文件成功")
	}
}

func TestDeploy(t *testing.T) {
	amount := new(big.Int).SetUint64(uint64(100000))
	decimal := uint8(8)
	content, err := ioutil.ReadFile(test.Erc20JsonLocation)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var unmarshalResponse ContractAbi
	err = json.Unmarshal(content, &unmarshalResponse)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	bytecode := unmarshalResponse.Bytecode
	contract, err := test.Web3.Thk.NewContract(unmarshalResponse.Abi)
	from := test.Web3.Thk.DefaultAddress
	nonce, err := test.Web3.Thk.GetNonce(from, test.Web3.Thk.DefaultChainId)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	transaction := util.Transaction{
		ChainId: test.Web3.Thk.DefaultChainId, FromChainId: test.Web3.Thk.DefaultChainId, ToChainId: test.Web3.Thk.DefaultChainId, From: from,
		To: "", Value: "0", Input: "", Nonce: strconv.Itoa(int(nonce)),
	}

	hash, err := contract.Deploy(transaction, bytecode, test.Web3.Thk.DefaultPrivateKey, Symbol, Name, decimal, amount)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("contract hash:", hash)
	time.Sleep(time.Second * 10)
	receipt, err := test.Web3.Thk.GetTransactionByHash(test.Web3.Thk.DefaultChainId, hash)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("contract addr:", receipt.ContractAddress)
	transaction = util.Transaction{
		ChainId: test.Web3.Thk.DefaultChainId, FromChainId: test.Web3.Thk.DefaultChainId, ToChainId: test.Web3.Thk.DefaultChainId, From: from,
		To: receipt.ContractAddress, Value: "0", Input: "", Nonce: strconv.Itoa(int(nonce)),
	}
	result, err := contract.Call(transaction, "symbol", test.Web3.Thk.DefaultPrivateKey)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("result:", result)

	var str string
	err = contract.Parse(result.Out, "symbol", &str)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var token Token
	token.Name = Name
	token.Symbol = Symbol
	token.Total = 100000
	token.ContractAddress = receipt.ContractAddress
	token.ABI = unmarshalResponse.Abi
	token.Icon = "icon"
	token.Website = "www.thinkey.com"
	token.Introduction = Des
	token.Date = time.Now().Format("2006-01-02")
	token.ChainId = "2"
	token.Decimal = 8
	PostInfo(token)
	//PostFile(token, "../resources/abc.sol")
}

func PostInfo(token Token) {
	stringUrl := "http://ext.thinkey.xyz/v1/wallet/token/tokeninfo/"
	bodybuf := bytes.NewBufferString("")
	bodywriter := multipart.NewWriter(bodybuf)
	bodywriter.SetBoundary("Pp7Ye2EeWaFDdAY")
	err := bodywriter.WriteField("Name", token.Name)
	err = bodywriter.WriteField("Symbol", token.Symbol)
	err = bodywriter.WriteField("Total", fmt.Sprintf("%f", token.Total))
	err = bodywriter.WriteField("ContractAddress", token.ContractAddress)
	err = bodywriter.WriteField("ABI", token.ABI)
	err = bodywriter.WriteField("Icon", token.Icon)
	err = bodywriter.WriteField("Website", token.Website)
	err = bodywriter.WriteField("Introduction", token.Introduction)
	err = bodywriter.WriteField("Date", token.Date)
	err = bodywriter.WriteField("ChainId", token.ChainId)
	err = bodywriter.WriteField("Decimal", fmt.Sprintf("%d", token.Decimal))
	bodywriter.Close()
	reqreader := io.MultiReader(bodybuf)
	resp, err := http.Post(stringUrl,
		bodywriter.FormDataContentType(),
		reqreader)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取回应消息异常")
	}
	str := string(body)
	fmt.Println("发送回应数据:", str)
}

func TestRun(t *testing.T) {
	var token Token
	token.Name = Name
	token.Symbol = Symbol
	token.Total = 100000
	token.ContractAddress = ""
	token.ABI = ""
	token.Icon = "icon"
	token.Website = ""
	PostFile(token, "../resources/abc.sol")
}
func PostFile(token Token, name string) {
	filebytes, err := ioutil.ReadFile(name)
	stringUrl := "http://192.168.1.164:8201/v1/thkadmin/wallet/token/"
	bodybuf := bytes.NewBufferString("")
	bodywriter := multipart.NewWriter(bodybuf)
	bodywriter.SetBoundary("Pp7Ye2EeWaFDdAY")
	err = bodywriter.WriteField("Name", token.Name)
	err = bodywriter.WriteField("Symbol", token.Symbol)
	err = bodywriter.WriteField("Total", fmt.Sprintf("%f", token.Total))
	err = bodywriter.WriteField("ContractAddress", token.ContractAddress)
	err = bodywriter.WriteField("ABI", token.ABI)
	err = bodywriter.WriteField("Icon", token.Icon)
	err = bodywriter.WriteField("Website", token.Website)
	filename := path.Base(name)
	_, err = bodywriter.CreateFormFile(token.Icon, filename)
	if err != nil {
		fmt.Printf("创建FormFile1文件信息异常！")
	}
	bodybuf.Write(filebytes)
	bodywriter.Close()
	//application/json
	//multipart/form-data

	reqreader := io.MultiReader(bodybuf)

	//req, err := http.NewRequest("POST", stringUrl, reqreader)
	//if err != nil {
	//	fmt.Printf("站点相机上传图片，创建上次请求异常！异常信息")
	//
	//}
	//req.Header.Set("Connection", "close")
	//req.Header.Set("Pragma", "no-cache")
	//req.Header.Set("Content-Type", bodywriter.FormDataContentType())
	//req.ContentLength = int64(bodybuf.Len())
	//fmt.Printf("发送消息长度:")
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//defer resp.Body.Close()

	resp, err := http.Post(stringUrl,
		bodywriter.FormDataContentType(),
		reqreader)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取回应消息异常")
	}
	fmt.Println("发送回应数据:", string(body))

}

func GetContractJson(ctct map[string]interface{}) (ContractJson, ContractJson, error) {
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
