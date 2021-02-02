package thk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/stephenfire/go-rtl"
	"io"
	"math/big"
	"strconv"
	"web3.go/common"
	"web3.go/common/hexutil"
)

const (
	SystemContractAddressWithdraw = "0x0000000000000000000000000000000000020000"
	SystemContractAddressDeposit  = "0x0000000000000000000000000000000000030000"
	SystemContractAddressCancel   = "0x0000000000000000000000000000000000040000"
)

type (
	MyChainId uint32
	MyHeight  uint64

	Addresser interface {
		Address() common.Address
	}
)

type MyCashCheque struct {
	FromChain    MyChainId      `json:"FromChain"`
	FromAddress  common.Address `json:"FromAddr"`
	Nonce        uint64         `json:"Nonce"`
	ToChain      MyChainId      `json:"ToChain"`
	ToAddress    common.Address `json:"ToAddr"`
	ExpireHeight MyHeight       `json:"ExpireHeight"` // The overdue height refers to that when the height of the target chain exceeds (excluding) this value, the check cannot be withdrawn and can only be returned
	Amount       *big.Int       `json:"Amount"`
}

type CashCheque struct {
	ChainId      string `json:"chainId"`
	FromChainId  string `json:"fromChainId"`
	From         string `json:"from"`
	Nonce        string `json:"nonce"`
	ToChainId    string `json:"toChainId"`
	To           string `json:"to"`
	ExpireHeight string `json:"expireheight"`
	Value        string `json:"value"`
}

func (c *CashCheque) Encode() (string, error) {
	fromBytes, err := hexutil.Decode(c.From)
	if err != nil {
		return "", err
	}
	toBytes, err := hexutil.Decode(c.To)
	if err != nil {
		return "", err
	}
	bigValue, ok := big.NewInt(0).SetString(c.Value, 10)
	if !ok {
		return "", fmt.Errorf("cheque value error :%v", c.Value)
	}
	fromChainId, err := strconv.Atoi(c.FromChainId)
	if err != nil {
		return "", err
	}
	toChainId, err := strconv.Atoi(c.ToChainId)
	if err != nil {
		return "", err
	}
	nonce, err := strconv.Atoi(c.Nonce)
	if err != nil {
		return "", err
	}
	expireHeight, err := strconv.Atoi(c.ExpireHeight)
	if err != nil {
		return "", err
	}
	cashCheque := &MyCashCheque{
		FromChain:    MyChainId(fromChainId),
		FromAddress:  common.BytesToAddress(fromBytes),
		Nonce:        uint64(nonce),
		ToChain:      MyChainId(toChainId),
		ToAddress:    common.BytesToAddress(toBytes),
		ExpireHeight: MyHeight(expireHeight),
		Amount:       bigValue,
	}
	chequeBytes, err := rtl.Marshal(cashCheque)
	return hexutil.Encode(chequeBytes), nil
}

func (c *CashCheque) Decode(input string) error {
	b, err := hexutil.Decode(input)
	if err != nil {
		return err
	}
	var cash MyCashCheque
	err = cash.Deserialization(bytes.NewReader(b))
	if err != nil {
		return err
	}
	c.ChainId = strconv.Itoa(int(cash.FromChain))
	c.FromChainId = strconv.Itoa(int(cash.FromChain))
	c.From = hexutil.Encode(cash.FromAddress[:])
	c.Nonce = strconv.FormatInt(int64(cash.Nonce), 10)
	c.ToChainId = strconv.Itoa(int(cash.ToChain))
	c.To = hexutil.Encode(cash.ToAddress[:])
	c.ExpireHeight = strconv.FormatInt(int64(cash.ExpireHeight), 10)
	c.Value = cash.Amount.String()
	return nil
}

// 4bytes FromChain + 20bytes FromAddress + 8bytes Nonce + 4bytes ToChain + 20bytes ToAddress +
// 8bytes ExpireHeight + 1bytes len(Amount.Bytes()) + Amount.Bytes()
// They are all bigendian
func (c *MyCashCheque) Serialization(w io.Writer) error {
	buf4 := make([]byte, 4)
	buf8 := make([]byte, 8)

	// 4bytes FromChain
	binary.BigEndian.PutUint32(buf4, uint32(c.FromChain))
	_, err := w.Write(buf4)
	if err != nil {
		return err
	}

	// 20bytes FromAddress
	_, err = w.Write(c.FromAddress[:])
	if err != nil {
		return err
	}

	// 8bytes Nonce
	binary.BigEndian.PutUint64(buf8, uint64(c.Nonce))
	_, err = w.Write(buf8)
	if err != nil {
		return err
	}

	// 4bytes ToChain
	binary.BigEndian.PutUint32(buf4, uint32(c.ToChain))
	_, err = w.Write(buf4)
	if err != nil {
		return err
	}

	// 20bytes ToAddress
	_, err = w.Write(c.ToAddress[:])
	if err != nil {
		return err
	}

	// 8bytes ExpireHeight
	binary.BigEndian.PutUint64(buf8, uint64(c.ExpireHeight))
	_, err = w.Write(buf8)
	if err != nil {
		return err
	}

	// 1bytes len(Amount.Bytes())
	buf4 = buf4[:1]
	var amountBytes []byte
	if c.Amount != nil {
		amountBytes = c.Amount.Bytes()
	}
	buf4[0] = byte(len(amountBytes))
	_, err = w.Write(buf4)
	if err != nil {
		return err
	}
	// Amount.Bytes()
	if buf4[0] > 0 {
		_, err = w.Write(amountBytes)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *MyCashCheque) Deserialization(r io.Reader) error {
	buf4 := make([]byte, 4)
	buf8 := make([]byte, 8)

	_, err := r.Read(buf4)
	if err != nil {
		return err
	}
	c.FromChain = MyChainId(binary.BigEndian.Uint32(buf4))

	_, err = r.Read(c.FromAddress[:])
	if err != nil {
		return err
	}

	_, err = r.Read(buf8)
	if err != nil {
		return err
	}
	c.Nonce = binary.BigEndian.Uint64(buf8)

	_, err = r.Read(buf4)
	if err != nil {
		return err
	}
	c.ToChain = MyChainId(binary.BigEndian.Uint32(buf4))

	_, err = r.Read(c.ToAddress[:])
	if err != nil {
		return err
	}

	_, err = r.Read(buf8)
	if err != nil {
		return err
	}
	c.ExpireHeight = MyHeight(binary.BigEndian.Uint64(buf8))

	buf4 = buf4[:1]
	_, err = r.Read(buf4)
	if err != nil {
		return err
	}
	length := int(buf4[0])

	if length > 0 {
		mbytes := make([]byte, length)
		_, err = r.Read(mbytes)
		if err != nil {
			return err
		}
		c.Amount = new(big.Int)
		c.Amount.SetBytes(mbytes)
	} else {
		c.Amount = big.NewInt(0)
	}

	return nil
}
