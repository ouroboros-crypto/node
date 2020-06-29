package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/magiconair/properties/assert"
	"github.com/ouroboros-crypto/node/x/coins"
	"testing"
)


// Tests the addToStructure method with default coin (ouro)
func TestAddToStructureDefault(t *testing.T) {
	ctx, _, keeper, _ := createTestInput(t, false)

	// Generate a few accounts
	_, addrs, _, _ := mock.CreateGenAccounts(2, sdk.NewCoins())

	owner := addrs[0]
	receiver := addrs[1]

	defaultCoin := coins.GetDefaultCoin()
	customCoin := coins.Coin{Symbol:"test"}

	// 100 ouro
	amount := sdk.NewIntWithDecimal(100, 6)

	added := keeper.AddToStructure(ctx, owner, receiver, amount, defaultCoin)

	// The receiver should be added to the owner structure
	assert.Equal(t, added, true)

	// Make sure we've got the right upper structure record
	upper := keeper.GetUpperStructure(ctx, receiver)

	assert.Equal(t, upper.Owner, owner)

	// Make sure the structure has everything right
	structure := keeper.GetStructure(ctx, owner, defaultCoin)

	assert.Equal(t, structure.Owner, owner)
	assert.Equal(t, structure.Balance, amount)
	assert.Equal(t, structure.Followers, sdk.NewInt(1))
	assert.Equal(t, structure.MaxLevel, sdk.NewInt(1))

	// Make sure any other custom coin will have 0 balance but still the same structure
	structure = keeper.GetStructure(ctx, owner, customCoin)

	assert.Equal(t, structure.Owner, owner)
	assert.Equal(t, structure.Balance, sdk.NewInt(0))
	assert.Equal(t, structure.Followers, sdk.NewInt(1))
	assert.Equal(t, structure.MaxLevel, sdk.NewInt(1))

	// It should not be saved again, even with another coint
	addedAgain := keeper.AddToStructure(ctx, owner, receiver, amount, customCoin)

	assert.Equal(t, addedAgain, false)
}

// Tests the addToStructure method with default coin (ouro)
func TestAddToStructureCustom(t *testing.T) {
	ctx, _, keeper, _ := createTestInput(t, false)

	// Generate a few accounts
	_, addrs, _, _ := mock.CreateGenAccounts(2, sdk.NewCoins())

	owner := addrs[0]
	receiver := addrs[1]

	customCoin := coins.Coin{Symbol:"test"}
	defaultCoin := coins.GetDefaultCoin()

	// 100 ouro
	amount := sdk.NewIntWithDecimal(100, 6)

	added := keeper.AddToStructure(ctx, owner, receiver, amount, customCoin)

	// The receiver should be added to the owner structure
	assert.Equal(t, added, true)

	// Make sure we've got the right upper structure record
	upper := keeper.GetUpperStructure(ctx, receiver)

	assert.Equal(t, upper.Owner, owner)

	// Make sure the structure has everything right
	structure := keeper.GetStructure(ctx, owner, customCoin)

	assert.Equal(t, structure.Owner, owner)
	assert.Equal(t, structure.Balance, amount)
	assert.Equal(t, structure.Followers, sdk.NewInt(1))
	assert.Equal(t, structure.MaxLevel, sdk.NewInt(1))

	// Make sure the default coin will have 0 balance but still the same structure
	structure = keeper.GetStructure(ctx, owner, defaultCoin)

	assert.Equal(t, structure.Owner, owner)
	assert.Equal(t, structure.Balance, sdk.NewInt(0))
	assert.Equal(t, structure.Followers, sdk.NewInt(1))
	assert.Equal(t, structure.MaxLevel, sdk.NewInt(1))

	// It should not be saved again, even with default coin
	addedAgain := keeper.AddToStructure(ctx, owner, receiver, amount, defaultCoin)

	assert.Equal(t, addedAgain, false)
}

// Tests the addToStructure method with default coin (ouro)
func TestAddToMaxStructure(t *testing.T) {
	ctx, _, keeper, _ := createTestInput(t, false)

	// Generate a few accounts
	_, addrs, _, _ := mock.CreateGenAccounts(152, sdk.NewCoins())

	// 10 ouro
	amount := sdk.NewIntWithDecimal(10, 6)

	coin := coins.GetDefaultCoin()

	// its depth will be ~150 levels (but will be reduced to 100 for the first account in line)
	i := 150

	// 150 -> 149, 149 -> 148, ...
	for i > 0 {
		assert.Equal(t, keeper.AddToStructure(ctx, addrs[i], addrs[i-1], amount, coin), true)

		i -= 1
	}

	// Make sure the first account doesn't have more than 100 followers & levels
	lastStructure := keeper.GetStructure(ctx, addrs[150], coin)

	assert.Equal(t, lastStructure.MaxLevel, sdk.NewInt(100))
	assert.Equal(t, lastStructure.Followers, sdk.NewInt(100))
	assert.Equal(t, lastStructure.Balance, sdk.NewInt(0)) // Since we've sent the money down the line

	// To make sure we'll get this balance removed
	lastStructure.Balance = sdk.NewIntWithDecimal(30, 6)

	keeper.SetStructure(ctx, lastStructure, coin)

	currentBalance := lastStructure.Balance

	// Make sure that adding follower to ~101 account won't be added to the first account
	assert.Equal(t, keeper.AddToStructure(ctx, addrs[50], addrs[151], amount, coin), true)

	lastStructure = keeper.GetStructure(ctx, addrs[150], coin)

	assert.Equal(t, lastStructure.MaxLevel.Equal(sdk.NewInt(100)), true)
	assert.Equal(t, lastStructure.Followers.Equal(sdk.NewInt(100)), true)
	assert.Equal(t, lastStructure.Balance, currentBalance.Sub(amount))
}

// Testing the increaseStructureBalance method with the default coin
func TestIncreaseStructureBalanceDefault(t *testing.T) {
	ctx, _, keeper, _ := createTestInput(t, false)

	// Generate a few accounts
	_, addrs, _, _ := mock.CreateGenAccounts(6, sdk.NewCoins())

	//customCoin := coins.Coin{Symbol:"test"}
	defaultCoin := coins.GetDefaultCoin()


	// 0 OURO
	coinsAmount := sdk.NewInt(0)

	// 5 levels
	i := 5

	// First let's create 5 levels depth structure - 5 -> 4, 4 -> 3, 3 -> 2, 2 -> 1, 1 -> 0
	for i > 0 {
		firstAccount, secondAccount := addrs[i], addrs[i-1]

		assert.Equal(t, keeper.AddToStructure(ctx, firstAccount, secondAccount, coinsAmount, defaultCoin), true)

		i -= 1
	}

	// The latest account in the structure gets 100 OURO
	coinsAmount = sdk.NewIntWithDecimal(100, 6)
	keeper.IncreaseStructureBalance(ctx, addrs[0], coinsAmount, defaultCoin)

	assert.Equal(t, keeper.GetStructure(ctx, addrs[0], defaultCoin).Balance.IsZero(), true)

	i = 5

	// Checking the balances
	for i > 1 {
		assert.Equal(t, keeper.GetStructure(ctx, addrs[i], defaultCoin).Balance.Equal(coinsAmount), true)

		i -= 1
	}
}

// Testing the increaseStructureBalance method with the custom coin
func TestIncreaseStructureBalanceCustom(t *testing.T) {
	ctx, _, keeper, _ := createTestInput(t, false)

	// Generate a few accounts
	_, addrs, _, _ := mock.CreateGenAccounts(6, sdk.NewCoins())

	coin := coins.Coin{Symbol:"test"}


	// 0test
	coinsAmount := sdk.NewInt(0)

	// 5 levels
	i := 5

	// First let's create 5 levels depth structure - 5 -> 4, 4 -> 3, 3 -> 2, 2 -> 1, 1 -> 0
	for i > 0 {
		firstAccount, secondAccount := addrs[i], addrs[i-1]

		assert.Equal(t, keeper.AddToStructure(ctx, firstAccount, secondAccount, coinsAmount, coin), true)

		i -= 1
	}

	// The latest account in the structure gets 100test
	coinsAmount = sdk.NewIntWithDecimal(100, 6)
	keeper.IncreaseStructureBalance(ctx, addrs[0], coinsAmount, coin)

	assert.Equal(t, keeper.GetStructure(ctx, addrs[0], coin).Balance.IsZero(), true)

	i = 5

	// Checking the balances
	for i > 1 {
		assert.Equal(t, keeper.GetStructure(ctx, addrs[i], coin).Balance.Equal(coinsAmount), true)

		i -= 1
	}
}

// Testing the increaseStructureBalance method with the default coin
func TestDecreaseStructureBalanceDefault(t *testing.T) {
	ctx, _, keeper, _ := createTestInput(t, false)

	// Generate a few accounts
	_, addrs, _, _ := mock.CreateGenAccounts(6, sdk.NewCoins())

	//customCoin := coins.Coin{Symbol:"test"}
	coin := coins.GetDefaultCoin()


	// 0 OURO
	coinsAmount := sdk.NewInt(0)

	// 5 levels
	i := 5

	// First let's create 5 levels depth structure - 5 -> 4, 4 -> 3, 3 -> 2, 2 -> 1, 1 -> 0
	for i > 0 {
		firstAccount, secondAccount := addrs[i], addrs[i-1]

		assert.Equal(t, keeper.AddToStructure(ctx, firstAccount, secondAccount, coinsAmount, coin), true)

		i -= 1
	}

	// The latest account in the structure gets 100 OURO
	coinsAmount = sdk.NewIntWithDecimal(100, 6)
	keeper.IncreaseStructureBalance(ctx, addrs[0], coinsAmount, coin)

	// And then transfers 40 OURO to another structure, now it has just 60 OURO in the structure
	keeper.DecreaseStructureBalance(ctx, addrs[0], sdk.NewIntWithDecimal(40, 6), coin)

	assert.Equal(t, keeper.GetStructure(ctx, addrs[0], coin).Balance.IsZero(), true)

	i = 5

	// Checking the balances
	for i > 1 {
		assert.Equal(t, keeper.GetStructure(ctx, addrs[i], coin).Balance.Equal(sdk.NewIntWithDecimal(60, 6)), true)

		i -= 1
	}
}


// Testing the increaseStructureBalance method with the default coin
func TestDecreaseStructureBalanceCustom(t *testing.T) {
	ctx, _, keeper, _ := createTestInput(t, false)

	// Generate a few accounts
	_, addrs, _, _ := mock.CreateGenAccounts(6, sdk.NewCoins())

	coin := coins.Coin{Symbol:"test"}

	// 0 OURO
	coinsAmount := sdk.NewInt(0)

	// 5 levels
	i := 5

	// First let's create 5 levels depth structure - 5 -> 4, 4 -> 3, 3 -> 2, 2 -> 1, 1 -> 0
	for i > 0 {
		firstAccount, secondAccount := addrs[i], addrs[i-1]

		assert.Equal(t, keeper.AddToStructure(ctx, firstAccount, secondAccount, coinsAmount, coin), true)

		i -= 1
	}

	// The latest account in the structure gets 100 OURO
	coinsAmount = sdk.NewIntWithDecimal(100, 6)
	keeper.IncreaseStructureBalance(ctx, addrs[0], coinsAmount, coin)

	// And then transfers 40 OURO to another structure, now it has just 60 OURO in the structure
	keeper.DecreaseStructureBalance(ctx, addrs[0], sdk.NewIntWithDecimal(40, 6), coin)

	assert.Equal(t, keeper.GetStructure(ctx, addrs[0], coin).Balance.IsZero(), true)

	i = 5

	// Checking the balances
	for i > 1 {
		assert.Equal(t, keeper.GetStructure(ctx, addrs[i], coin).Balance.Equal(sdk.NewIntWithDecimal(60, 6)), true)

		i -= 1
	}
}