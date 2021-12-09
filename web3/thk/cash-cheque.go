package thk

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ThinkiumGroup/web3.go/common"
	"github.com/ThinkiumGroup/web3.go/common/hexutil"
	"github.com/ThinkiumGroup/web3.go/common/math"
	"github.com/stephenfire/go-rtl"
	"io"
	"math/big"
	"strconv"
)

const (
	SystemContractAddressWithdraw = "0x0000000000000000000000000000000000020000"
	SystemContractAddressDeposit  = "0x0000000000000000000000000000000000030000"
	SystemContractAddressCancel   = "0x0000000000000000000000000000000000040000"
)

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
	cashCheque := &CashCheck{
		FromChain:    common.ChainId(fromChainId),
		FromAddress:  common.BytesToAddress(fromBytes),
		Nonce:        uint64(nonce),
		ToChain:      common.ChainId(toChainId),
		ToAddress:    common.BytesToAddress(toBytes),
		ExpireHeight: common.Height(expireHeight),
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
	var cash CashCheck
	if len(input) < 200 { // TODO
		_, err = cash.Deserialization(bytes.NewReader(b))
		if err != nil {
			return err
		}
	} else {
		request := new(CashRequest)
		b, err := hexutil.Decode(input)
		if err != nil {
			return err
		}
		if err = rtl.Unmarshal(b, request); err != nil {
			return err
		}
		cash = *request.Check
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

type CashCheck struct {
	ParentChain  common.ChainId `json:"ParentChain"`
	IsShard      bool           `json:"IsShard"`
	FromChain    common.ChainId `json:"FromChain"`
	FromAddress  common.Address `json:"FromAddr"`
	Nonce        uint64         `json:"Nonce"`
	ToChain      common.ChainId `json:"ToChain"`
	ToAddress    common.Address `json:"ToAddr"`
	ExpireHeight common.Height  `json:"ExpireHeight"`
	UserLocal    bool           `json:"UseLocal"`
	Amount       *big.Int       `json:"Amount"`
	CurrencyID   common.ChainId `json:"CoinID"`
}

func (c *CashCheck) String() string {
	return fmt.Sprintf("Check{ParentChain:%d IsShard:%t From:[%d,%x] Nonce:%d To:[%d,%x]"+" Expire:%d Local:%t Amount:%s CoinID:%d}",
		c.ParentChain, c.IsShard, c.FromChain, c.FromAddress[:], c.Nonce, c.ToChain, c.ToAddress[:], c.ExpireHeight, c.UserLocal, math.BigIntForPrint(c.Amount), c.CurrencyID)
}

func (c *CashCheck) serialPrefix() []byte {
	var version byte = 0x0
	if c.ParentChain == 0 && c.IsShard == false && c.CurrencyID == 0 {
		if c.UserLocal == false {
			return nil
		} else {
			version = 0x0
			panic("wrong data: UseLocal==true without CurrencyId")
		}
	} else {
		version = 0x1
		buf := make([]byte, 13)
		binary.BigEndian.PutUint32(buf[:4], uint32(common.ReservedMaxChainID))
		buf[4] = version
		if c.UserLocal {
			buf[5] = 0x1
		} else {
			buf[5] = 0x0
		}
		binary.BigEndian.PutUint32(buf[6:10], uint32(c.ParentChain))
		if c.IsShard {
			buf[10] = 0x1
		} else {
			buf[10] = 0x0
		}
		binary.BigEndian.PutUint16(buf[11:13], uint16(c.CurrencyID))
		return buf
	}
}

func (c *CashCheck) Serialization(w io.Writer) error {
	buf4 := make([]byte, 4)
	buf8 := make([]byte, 8)

	var err error
	prefix := c.serialPrefix()
	if len(prefix) > 0 {
		_, err = w.Write(prefix)
		if err != nil {
			return err
		}
	}

	binary.BigEndian.PutUint32(buf4, uint32(c.FromChain))
	_, err = w.Write(buf4)
	if err != nil {
		return err
	}

	_, err = w.Write(c.FromAddress[:])
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint64(buf8, uint64(c.Nonce))
	_, err = w.Write(buf8)
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint32(buf4, uint32(c.ToChain))
	_, err = w.Write(buf4)
	if err != nil {
		return err
	}

	_, err = w.Write(c.ToAddress[:])
	if err != nil {
		return err
	}

	binary.BigEndian.PutUint64(buf8, uint64(c.ExpireHeight))
	_, err = w.Write(buf8)
	if err != nil {
		return err
	}

	var mbytes []byte
	if c.Amount != nil {
		mbytes = c.Amount.Bytes()
	}
	err = writeByteSlice(w, 1, mbytes)
	if err != nil {
		return err
	}

	return nil
}

func (c *CashCheck) Deserialization(r io.Reader) (shouldBeNil bool, err error) {
	buf4 := make([]byte, 4)
	buf8 := make([]byte, 8)

	_, err = io.ReadFull(r, buf4)
	if err != nil {
		return
	}
	first := common.ChainId(binary.BigEndian.Uint32(buf4))
	if first.IsNil() {
		_, err = io.ReadFull(r, buf4[:1])
		if err != nil {
			return
		}

		switch buf4[0] {
		case 0x0:
			c.UserLocal = true
			c.ParentChain = 0
			c.IsShard = false
			c.CurrencyID = 0
		case 0x1:
			_, err = io.ReadFull(r, buf8)
			if err != nil {
				return
			}
			if buf8[0] == 0x0 {
				c.UserLocal = false
			} else {
				c.UserLocal = true
			}
			c.ParentChain = common.ChainId(binary.BigEndian.Uint32(buf8[1:5]))
			if buf8[5] == 0x0 {
				c.IsShard = false
			} else {
				c.IsShard = true
			}
			c.CurrencyID = common.ChainId(binary.BigEndian.Uint16(buf8[6:8]))
		default:
			err = fmt.Errorf("unknown version of check %x", buf4[0])
			return
		}

		_, err = io.ReadFull(r, buf4)
		if err != nil {
			return
		}
		c.FromChain = common.ChainId(binary.BigEndian.Uint32(buf4))
	} else {
		c.FromChain = first
	}

	_, err = io.ReadFull(r, c.FromAddress[:])
	if err != nil {
		return
	}

	_, err = io.ReadFull(r, buf8)
	if err != nil {
		return
	}
	c.Nonce = binary.BigEndian.Uint64(buf8)

	_, err = io.ReadFull(r, buf4)
	if err != nil {
		return
	}
	c.ToChain = common.ChainId(binary.BigEndian.Uint32(buf4))

	_, err = io.ReadFull(r, c.ToAddress[:])
	if err != nil {
		return
	}

	_, err = io.ReadFull(r, buf8)
	if err != nil {
		return
	}
	c.ExpireHeight = common.Height(binary.BigEndian.Uint64(buf8))

	bs, err := readByteSlice(r, 1)
	if err != nil {
		return false, err
	}
	if len(bs) > 0 {
		c.Amount = new(big.Int).SetBytes(bs)
	} else {
		c.Amount = big.NewInt(0)
	}

	return false, nil
}

type CashRequest struct {
	Check *CashCheck `json:"check"`
}

func writeByteSlice(w io.Writer, uintType int, bs []byte) error {
	n := len(bs)
	var buf []byte
	switch uintType {
	case 1:
		if n > 0xFF {
			return errors.New("length is too big")
		}
		buf = make([]byte, 1)
		buf[0] = byte(n)
	case 2:
		if n > 0xFFFF {
			return errors.New("length is too big")
		}
		buf = make([]byte, 2)
		binary.BigEndian.PutUint16(buf, uint16(n))
	case 4:
		if n > 0xFFFFFFFF {
			return errors.New("length is too big")
		}
		buf = make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(n))
	case 8:
		buf = make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(n))
	default:
		return errors.New("unknown type")
	}
	_, err := w.Write(buf)
	if err != nil {
		return err
	}
	if n > 0 {
		_, err = w.Write(bs)
		if err != nil {
			return err
		}
	}
	return nil
}

func readByteSlice(r io.Reader, uintType int) (bs []byte, err error) {
	var buf []byte
	switch uintType {
	case 1:
		buf = make([]byte, 1)
	case 2:
		buf = make([]byte, 2)
	case 4:
		buf = make([]byte, 4)
	case 8:
		buf = make([]byte, 8)
	default:
		return nil, errors.New("unknown type")
	}
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	var n uint64
	switch uintType {
	case 1:
		n = uint64(buf[0])
	case 2:
		n = uint64(binary.BigEndian.Uint16(buf))
	case 4:
		n = uint64(binary.BigEndian.Uint32(buf))
	case 8:
		n = uint64(binary.BigEndian.Uint64(buf))
	default:
		return nil, errors.New("unknown type")
	}
	if n > 0 {
		bs = make([]byte, n)
		_, err = io.ReadFull(r, bs)
		if err != nil {
			return nil, err
		}
		return bs, nil
	}
	return nil, nil
}
