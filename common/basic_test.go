package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	errorCheckSumAddress = "0x08d04fc4513e27854ac19f16B9b8cB8D564e2c68"
	rightCheckSumAddress = "0x08D04fc4513e27854ac19f16B9b8cB8D564e2c68"
	hex1                 = "8D04fc4513e27854ac19f16B9b8cB8D564e2c68"
	strictAddress        = "0x08d04fc4513e27854ac19f16b9b8cb8d564e2c68"
)

func TestIsChecksumAddress(t *testing.T) {
	assert.True(t, IsChecksumAddress(rightCheckSumAddress))
	assert.False(t, IsChecksumAddress(errorCheckSumAddress))
}

func TestIsStrictAddress(t *testing.T) {
	assert.True(t, IsStrictAddress(strictAddress))
	assert.False(t, IsStrictAddress(rightCheckSumAddress))
}

func TestToChecksumAddress(t *testing.T) {
	assert.Equal(t, rightCheckSumAddress, ToChecksumAddress(strictAddress))
}

func TestToStrictAddress(t *testing.T) {
	assert.Equal(t, strictAddress, ToStrictAddress(hex1))
	assert.Equal(t, strictAddress, ToStrictAddress(rightCheckSumAddress))
	assert.Equal(t, strictAddress, ToStrictAddress(errorCheckSumAddress))
}
