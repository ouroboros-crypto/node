package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ouroboros-crypto/node/x/posmining/types"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoinsPerTime(t *testing.T) {
	balance := sdk.NewIntWithDecimal(1000, 6) // 1000 ouro
	zeroInt := sdk.NewInt(0)

	perTime := types.NewCoinsPerTime(balance, sdk.NewInt(6), zeroInt, zeroInt, zeroInt)

	assert.Equal(t, perTime.Day, sdk.NewIntWithDecimal(6, 5)) // 0.6 per day
	assert.Equal(t, perTime.Hour, sdk.NewIntWithDecimal(25, 3)) // 0.025 per hour
	assert.Equal(t, perTime.Minute, sdk.NewInt(416)) // 0.000416 per minute
}

func TestTimeDifference(t *testing.T) {
	seconds := sdk.NewInt(86400 + 3600 + 60 + 10) // 1 day, 1 hour, 1 minute and 10 seconds

	difference := types.NewTimeDifference(seconds)

	assert.Equal(t, difference.Days, sdk.NewInt(1))
	assert.Equal(t, difference.Hours, sdk.NewInt(1))
	assert.Equal(t, difference.Minutes, sdk.NewInt(1))
	assert.Equal(t, difference.Seconds, sdk.NewInt(10))
}
