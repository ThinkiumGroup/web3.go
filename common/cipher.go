package common

import (
	"encoding/hex"
	"github.com/ThinkiumGroup/go-cipher"
	"github.com/ThinkiumGroup/web3.go/common/hexutil"
)

var (
	Cipher = cipher.NewCipher(cipher.SECP256K1SHA3)
)

func HexToPrivateKey(h string) (cipher.ECCPrivateKey, error) {
	h = CleanHexPrefix(h)
	bs, err := hex.DecodeString(h)
	if err != nil {
		return nil, err
	}
	return Cipher.BytesToPriv(bs)
}

func Hash256(s string) []byte {
	return SystemHash256([]byte(s))
}

func SystemHash256(in ...[]byte) []byte {
	hasher := Cipher.Hasher()
	for _, b := range in {
		hasher.Write(b)
	}
	return hasher.Sum(nil)
}

func Sign(s, privateKey string) (string, error) {
	hash := Hash256(s)
	key, err := HexToPrivateKey(privateKey)
	if err != nil {
		return "", err
	}
	sig, err := Cipher.Sign(Cipher.PrivToBytes(key), hash)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(sig), nil
}
