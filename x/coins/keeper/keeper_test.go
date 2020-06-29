package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/magiconair/properties/assert"
	"testing"
)

// Tests update creation price
func TestKeeper_UpdateCreationPrice(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false)

	// 0.02
	assert.Equal(t, keeper.UpdateCreationPrice(ctx, sdk.NewInt(2)), sdk.NewIntWithDecimal(2500, 6))

	creationPrice, _ := keeper.GetCreationPrice(ctx)

	assert.Equal(t, creationPrice.Price, sdk.NewIntWithDecimal(2500, 6))

	// 0.4 usd
	assert.Equal(t, keeper.UpdateCreationPrice(ctx, sdk.NewInt(40)), sdk.NewIntWithDecimal(125, 6))

	creationPrice, _ = keeper.GetCreationPrice(ctx)

	assert.Equal(t, creationPrice.Price, sdk.NewIntWithDecimal(125, 6))
}
