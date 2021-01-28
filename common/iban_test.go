package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	iBanAddress      = "TH93112A8NT0RIRPIXR7IXJBB1U0ZBXEB1K"
	checkSumAddress  = "0x08D04fc4513e27854ac19f16B9b8cB8D564e2c68"
	errorIBanAddress = "TH9312A8NT0RIRPIXR7IXJBB1U0ZBXEB1K"
)

func TestIsDirect(t *testing.T) {
	assert.True(t, IsDirect(iBanAddress))
}

func TestIsValid(t *testing.T) {
	assert.True(t, IsValid(iBanAddress))
	assert.False(t, IsValid(errorIBanAddress))
}

func TestToAddress(t *testing.T) {
	assert.Equal(t, checkSumAddress, ToAddress(iBanAddress))
}

func TestToIBan(t *testing.T) {
	assert.Equal(t, iBanAddress, ToIBan(checkSumAddress))
}
