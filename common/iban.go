package common

import (
	"github.com/ThinkiumGroup/web3.go/common/hexutil"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

//Prepare an IBAN for mod 97 computation by moving the first 4 chars to the end and transforming the letters to
//numbers (A = 10, B = 11, ..., Z = 35), as specified in ISO13616.
func iso13616Prepare(iBan string) (s string) {
	iBan = strings.ToUpper(iBan)
	iBan = iBan[4:] + iBan[0:4]
	for i := 0; i < len(iBan); i++ {
		num := iBan[i]
		if num >= 'A' && num <= 'Z' {
			// A = 10, B = 11, ... Z = 35
			s += strconv.Itoa(int(num - 'A' + 10))
		} else {
			s += string(iBan[i])
		}
	}
	return s
}
func mod9710(iBan string) int {
	var remainder, block = iBan, ""
	for {
		if len(remainder) <= 2 {
			break
		}

		block = remainder[0:int(math.Min(9, float64(len(remainder))))]
		num, _ := strconv.Atoi(block)
		remainder = strconv.Itoa(num%97) + remainder[len(block):]
	}
	num, _ := strconv.Atoi(remainder)
	return num % 97
}
func fromBBan(bBan string) string {
	countryCode := "TH"
	remainder := mod9710(iso13616Prepare(countryCode + "00" + bBan))
	checkDigit := strconv.Itoa(98 - remainder)
	return countryCode + padLeft(checkDigit, 2) + bBan
}
func fromAddress(address string) string {
	num, ok := new(big.Int).SetString(strings.ToLower(CleanHexPrefix(address)), 16)
	if !ok {
		return ""
	}
	padded := padLeft(num.Text(36), 30)
	return fromBBan(strings.ToUpper(padded))
}
func ToIBan(address string) string {
	return fromAddress(address)
}
func ToAddress(iBan string) string {
	if IsDirect(iBan) {
		num, ok := new(big.Int).SetString(iBan[4:], 36)
		if !ok {
			return ""
		}
		return ToChecksumAddress(padLeft(CleanHexPrefix(hexutil.EncodeBig(num)), 40))
	}
	return ""
}
func IsDirect(iBan string) bool {
	return len(iBan) == 34 || len(iBan) == 35
}
func IsValid(iBan string) bool {
	matched, err := regexp.MatchString("^TH[0-9]{2}[0-9A-Z]{30,31}$", iBan)
	if err != nil || !matched {
		return false
	}
	return mod9710(iso13616Prepare(iBan)) == 1
}
