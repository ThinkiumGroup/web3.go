package common

import (
	"encoding/hex"
	"web3.go/common/hexutil"
)

var (
	EmptyPlaceHolder = struct{}{}
)

type (
	Hash    [HashLength]byte
	Address [AddressLength]byte
)

func (h Hash) Bytes() []byte { return h[:] }

func (h Hash) Hex() string { return hexutil.Encode(h[:]) }

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

// BytesToHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

func (a Address) String() string {
	return hex.EncodeToString(a[:])
}

func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

func HexToAddress(s string) Address {
	if HasHexPrefix(s) {
		s = s[2:]
	}
	b, _ := hex.DecodeString(s)
	var a Address
	a.SetBytes(b)
	return a
}
