package common

import (
	"encoding/hex"
	"regexp"
	"strconv"
	"strings"
	"web3.go/common/hexutil"
)

var (
	EmptyPlaceHolder = struct{}{}
)

const pad = "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"

type (
	Hash    [HashLength]byte
	Address [AddressLength]byte
)

func padLeft(str string, chars int) string {
	if chars <= len(str) {
		return str
	}
	return pad[0:chars-len(str)] + str
}

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
	s = CleanHexPrefix(s)
	b, _ := hex.DecodeString(s)
	var a Address
	a.SetBytes(b)
	return a
}

func IsStrictAddress(address string) bool {
	matched, err := regexp.MatchString("^0x[0-9a-f]{40}$", address)
	if err != nil || !matched {
		return false
	}
	return true
}
func IsChecksumAddress(address string) bool {
	address = CleanHexPrefix(address)
	addressHash := strings.Split(CleanHexPrefix(hexutil.Encode(Hash256(strings.ToLower(address)))), "")
	for i := 0; i < 40; i++ {
		//the nth letter should be uppercase if the nth digit of casemap is 1
		num, _ := strconv.ParseInt(addressHash[i], 16, 64)
		c := address[i]
		if (num > 7 && c != strings.ToUpper(string(c))[0]) || (num <= 7 && c != strings.ToLower(string(c))[0]) {
			return false
		}
	}
	return true
}
func ToChecksumAddress(address string) (s string) {
	address = strings.ToLower(CleanHexPrefix(address))
	if !IsStrictAddress("0x" + address) {
		return ""
	}
	s = "0x"
	addressHash := strings.Split(CleanHexPrefix(hexutil.Encode(Hash256(address))), "")
	for i := 0; i < len(address); i++ {
		// If ith character is 9 to f then make it uppercase
		num, _ := strconv.ParseInt(addressHash[i], 16, 64)
		if num > 7 {
			s += strings.ToUpper(string(address[i]))
		} else {
			s += string(address[i])
		}
	}
	return s
}

func ToStrictAddress(address string) string {
	if IsStrictAddress(address) {
		return address
	}
	match, _ := regexp.MatchString("^[0-9a-f]{40}$", address)
	if match {
		return "0x" + address
	}

	return "0x" + padLeft(strings.ToLower(CleanHexPrefix(address)), 40)
}
