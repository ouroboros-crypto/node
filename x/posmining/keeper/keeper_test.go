package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/magiconair/properties/assert"
	"github.com/ouroboros-crypto/node/x/coins"
	"github.com/ouroboros-crypto/node/x/posmining/types"
	"testing"
	"time"
)

var (
	BasicQuo = sdk.NewInt(100)
)

func MinusDays(days int) time.Time {
	defaultDate := time.Date(2020, 1, 10, 00, 0, 0, 0, time.UTC)

	defaultDate = defaultDate.AddDate(0, 0, days*-1)

	return defaultDate
}

func GetTestAcct() sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32("ouro1tcwh5z59sg0rcsuavng80wpgkcegmpl3egcqlp")

	return addr
}

// Tests the get saving periods methods
func TestGetSavingPeriods(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false)

	currentTime := MinusDays(0)
	LastTransaction := MinusDays(29)

	// Less than 30 days
	periods := keeper.GetSavingPeriods(ctx.WithBlockTime(currentTime), types.Posmining{LastTransaction: LastTransaction, LastCharged: LastTransaction})

	assert.Equal(t, len(periods), 1)
	assert.Equal(t, periods[0].End, currentTime)
	assert.Equal(t, periods[0].Start, LastTransaction)
	assert.Equal(t, periods[0].SavingCoff, sdk.NewInt(0))
	assert.Equal(t, periods[0].CorrectionCoff, sdk.NewInt(0)) // It should be always 0
}

// Tests the get saving periods methods
func TestGetCorrectionPeriods(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false)

	currentTime := MinusDays(0)
	LastRegulation := MinusDays(5)
	LastCharge := MinusDays(4)

	regulation := types.Correction{StartDate: LastRegulation, OpeningPrice: sdk.NewInt(10), CorrectionCoff: sdk.NewInt(20)}

	keeper.SetCorrection(ctx, regulation)

	// Less than 30 days
	periods := keeper.GetCorrectionPeriods(ctx.WithBlockTime(currentTime), types.Posmining{LastTransaction: LastCharge, LastCharged: LastCharge})

	assert.Equal(t, len(periods), 1)
	assert.Equal(t, periods[0].End, currentTime)
	assert.Equal(t, periods[0].Start, LastCharge) // Since we have to make sure it's within our posmining period
	assert.Equal(t, periods[0].SavingCoff, sdk.NewInt(0))
	assert.Equal(t, periods[0].CorrectionCoff, regulation.CorrectionCoff)
}

// Testing the situation when both savings and regulations got just a single period
func TestGetPosminingGroup1(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false)

	coin := coins.GetDefaultCoin()

	currentTime := MinusDays(0)
	LastRegulation := MinusDays(10)
	LastCharge := MinusDays(4)
	LastTx := MinusDays(6)

	balance := sdk.NewIntWithDecimal(1000, 6)

	paraminig := types.Posmining{
		DailyPercent:    sdk.NewInt(16), // 0.16% per day or 1.6 if balance 1000
		StructureCoff:   sdk.NewInt(218),
		Paramined:       sdk.NewInt(0),
		LastTransaction: LastTx,
		LastCharged:     LastCharge,
	}

	// 0.0016 per day * 2.18 structure coff * 0.20 regulation coff
	shouldGetPosmined := balance.Mul(sdk.NewInt(16)).Quo(sdk.NewInt(10000)).
		Mul(sdk.NewInt(218)).Quo(sdk.NewInt(100)).
		Mul(sdk.NewInt(20)).Quo(sdk.NewInt(100)).Mul(sdk.NewInt(4))


	// 0.2
	regulation := types.Correction{StartDate: LastRegulation, OpeningPrice: sdk.NewInt(10), CorrectionCoff: sdk.NewInt(20)}

	keeper.SetCorrection(ctx, regulation)

	group := keeper.GetPosminingGroup(ctx.WithBlockTime(currentTime), paraminig, coin,balance)

	assert.Equal(t, len(group.Periods), 1)
	assert.Equal(t, group.Periods[0].SavingCoff, sdk.NewInt(0)) // since it's 70 since the latest outcoming tx
	assert.Equal(t, group.Paramined, shouldGetPosmined)
}

// Tests the situation when we got multiple savings periods but only one regulation
func TestGetPosminingGroup2(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false)

	coin := coins.GetDefaultCoin()

	//account := GetTestAcct()

	currentTime := MinusDays(0)
	LastRegulation := MinusDays(5)
	LastCharge := MinusDays(4)
	LastTx := MinusDays(70)

	balance := sdk.NewIntWithDecimal(1000, 6)

	paraminig := types.Posmining{
		DailyPercent:    sdk.NewInt(16), // 0.16% per day or 1.6 if balance 1000
		StructureCoff:   sdk.NewInt(218),
		Paramined:       sdk.NewInt(0),
		LastTransaction: LastTx,
		LastCharged:     LastCharge,
	}

	// 0.0016 per day * 2.18 structure coff * 1.51 saving coff * 0.20 regulation coff
	// 0.0016 * 2.18 * 1.51 * 0.2 = 0.001053, so 1000 * 0.0010 * 4
	shouldGetPosmined := balance.Mul(sdk.NewInt(16)).Quo(sdk.NewInt(10000)).
		Mul(sdk.NewInt(218)).Quo(sdk.NewInt(100)).Mul(sdk.NewInt(151)).Quo(sdk.NewInt(100)).
		Mul(sdk.NewInt(20)).Quo(sdk.NewInt(100)).Mul(sdk.NewInt(4))


	// 0.2
	regulation := types.Correction{StartDate: LastRegulation, OpeningPrice: sdk.NewInt(10), CorrectionCoff: sdk.NewInt(20)}

	keeper.SetCorrection(ctx, regulation)

	group := keeper.GetPosminingGroup(ctx.WithBlockTime(currentTime), paraminig, coin,balance)

	assert.Equal(t, len(group.Periods), 1)
	assert.Equal(t, group.Periods[0].SavingCoff, sdk.NewInt(151)) // since it's 70 since the latest outcoming tx
	assert.Equal(t, group.Paramined, shouldGetPosmined)
}


// Tests the situation when we got multiple regulation periods but only one saving
func TestGetPosminingGroup3(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false)

	coin := coins.GetDefaultCoin()

	//account := GetTestAcct()

	currentTime := MinusDays(0)
	LastRegulation := MinusDays(2)
	LastCharge := MinusDays(4)
	LastTx := MinusDays(4)

	balance := sdk.NewIntWithDecimal(1000, 6)

	paraminig := types.Posmining{
		DailyPercent:    sdk.NewInt(16), // 0.16% per day or 1.6 if balance 1000
		StructureCoff:   sdk.NewInt(218),
		Paramined:       sdk.NewInt(0),
		LastTransaction: LastTx,
		LastCharged:     LastCharge,
	}

	// 0.2
	previousRegulation := types.PreviousCorrection{
		StartDate:      LastCharge, // 2 day long
		EndDate:        LastRegulation,
		OpeningPrice:   sdk.NewInt(50),
		CorrectionCoff: sdk.NewInt(10), // 0.1
	}

	regulation := types.Correction{StartDate: LastRegulation, OpeningPrice: sdk.NewInt(10), CorrectionCoff: sdk.NewInt(20), PreviousCorrections: []types.PreviousCorrection{previousRegulation}}

	keeper.SetCorrection(ctx, regulation)

	basicPosmining := balance.Mul(paraminig.DailyPercent).Quo(sdk.NewInt(10000)).Mul(paraminig.StructureCoff).Quo(sdk.NewInt(100))

	firstHalf := basicPosmining.Mul(regulation.CorrectionCoff).Quo(BasicQuo).Mul(sdk.NewInt(2))
	secondHalf := basicPosmining.Mul(previousRegulation.CorrectionCoff).Quo(BasicQuo).Mul(sdk.NewInt(2))

	shouldGetPosmined := firstHalf.Add(secondHalf)

	group := keeper.GetPosminingGroup(ctx.WithBlockTime(currentTime), paraminig, coin,balance)

	assert.Equal(t, len(group.Periods), 2)
	assert.Equal(t, group.Periods[0].SavingCoff, sdk.NewInt(0))
	assert.Equal(t, group.Periods[0].CorrectionCoff, previousRegulation.CorrectionCoff)
	assert.Equal(t, group.Periods[1].SavingCoff, sdk.NewInt(0))
	assert.Equal(t, group.Periods[1].CorrectionCoff, regulation.CorrectionCoff)
	assert.Equal(t, group.Paramined, shouldGetPosmined)
}

// Tests the situation when we got multiple regulation and savings periods
func TestGetPosminingGroup4(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false)

	coin := coins.GetDefaultCoin()

	currentTime := MinusDays(0)
	LastRegulation := MinusDays(20)
	LastCharge := MinusDays(40)
	LastTx := MinusDays(70)

	balance := sdk.NewIntWithDecimal(1000, 6)

	paraminig := types.Posmining{
		DailyPercent:    sdk.NewInt(16), // 0.16% per day or 1.6 if balance 1000
		StructureCoff:   sdk.NewInt(218),
		Paramined:       sdk.NewInt(0),
		LastTransaction: LastTx,
		LastCharged:     LastCharge,
	}

	// 0.2
	previousRegulation := types.PreviousCorrection{
		StartDate:      LastCharge, // 2 day long
		EndDate:        LastRegulation,
		OpeningPrice:   sdk.NewInt(50),
		CorrectionCoff: sdk.NewInt(10), // 0.1
	}

	regulation := types.Correction{StartDate: LastRegulation, OpeningPrice: sdk.NewInt(10), CorrectionCoff: sdk.NewInt(20), PreviousCorrections: []types.PreviousCorrection{previousRegulation}}

	keeper.SetCorrection(ctx, regulation)

	basicPosmining := balance.Mul(paraminig.DailyPercent).Quo(sdk.NewInt(10000)).Mul(paraminig.StructureCoff).Quo(sdk.NewInt(100))

	// First period is 20 days long, it uses the previous regulation and saving coff 150
	firstPeriod := basicPosmining.Mul(previousRegulation.CorrectionCoff).Quo(BasicQuo).Mul(sdk.NewInt(150)).Quo(BasicQuo).Mul(sdk.NewInt(20))
	secondPeriod := basicPosmining.Mul(regulation.CorrectionCoff).Quo(BasicQuo).Mul(sdk.NewInt(150)).Quo(BasicQuo).Mul(sdk.NewInt(10))
	thirdPeriod := basicPosmining.Mul(regulation.CorrectionCoff).Quo(BasicQuo).Mul(sdk.NewInt(151)).Quo(BasicQuo).Mul(sdk.NewInt(10))

	shouldGetPosmined := firstPeriod.Add(secondPeriod).Add(thirdPeriod)

	group := keeper.GetPosminingGroup(ctx.WithBlockTime(currentTime), paraminig, coin,balance)

	assert.Equal(t, len(group.Periods), 3)

	assert.Equal(t, group.Periods[0].SavingCoff, sdk.NewInt(150))
	assert.Equal(t, group.Periods[0].CorrectionCoff, previousRegulation.CorrectionCoff)

	assert.Equal(t, group.Periods[1].SavingCoff, sdk.NewInt(150))
	assert.Equal(t, group.Periods[1].CorrectionCoff, regulation.CorrectionCoff)

	assert.Equal(t, group.Periods[2].SavingCoff, sdk.NewInt(151))
	assert.Equal(t, group.Periods[2].CorrectionCoff, regulation.CorrectionCoff)

	assert.Equal(t, group.Paramined, shouldGetPosmined)
}

// Tests the situation when we account just got his first tx, so there shouldn't be any periods
func TestGetPosminingGroup5(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false)

	coin := coins.GetDefaultCoin()

	currentTime := MinusDays(0)
	LastRegulation := MinusDays(20)
	LastCharge := MinusDays(40)

	balance := sdk.NewIntWithDecimal(1000, 6)

	paraminig := types.Posmining{
		DailyPercent:    sdk.NewInt(16), // 0.16% per day or 1.6 if balance 1000
		StructureCoff:   sdk.NewInt(218),
		Paramined:       sdk.NewInt(0),
		LastTransaction: currentTime,
		LastCharged:     currentTime,
	}

	// 0.2
	previousRegulation := types.PreviousCorrection{
		StartDate:      LastCharge, // 2 day long
		EndDate:        LastRegulation,
		OpeningPrice:   sdk.NewInt(50),
		CorrectionCoff: sdk.NewInt(10), // 0.1
	}

	regulation := types.Correction{StartDate: LastRegulation, OpeningPrice: sdk.NewInt(10), CorrectionCoff: sdk.NewInt(20), PreviousCorrections: []types.PreviousCorrection{previousRegulation}}

	keeper.SetCorrection(ctx, regulation)

	group := keeper.GetPosminingGroup(ctx.WithBlockTime(currentTime), paraminig, coin,balance)

	assert.Equal(t, len(group.Periods), 0)
}

// Tests the situation when we account just got his first tx and the time diff is in milliseconds
func TestGetPosminingGroup6(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false)

	coin := coins.GetDefaultCoin()

	LastRegulation := MinusDays(20)
	LastCharge := MinusDays(40)
	defaultTime := MinusDays(0)

	currentTime := time.Date(2020, 1, 10, 00, 0, 0, 10, time.UTC)

	balance := sdk.NewIntWithDecimal(1000, 6)

	paraminig := types.Posmining{
		DailyPercent:    sdk.NewInt(16), // 0.16% per day or 1.6 if balance 1000
		StructureCoff:   sdk.NewInt(218),
		Paramined:       sdk.NewInt(0),
		LastTransaction: defaultTime,
		LastCharged:     defaultTime,
	}

	// 0.2
	previousRegulation := types.PreviousCorrection{
		StartDate:      LastCharge, // 2 day long
		EndDate:        LastRegulation,
		OpeningPrice:   sdk.NewInt(50),
		CorrectionCoff: sdk.NewInt(10), // 0.1
	}

	regulation := types.Correction{StartDate: LastRegulation, OpeningPrice: sdk.NewInt(10), CorrectionCoff: sdk.NewInt(20), PreviousCorrections: []types.PreviousCorrection{previousRegulation}}

	keeper.SetCorrection(ctx, regulation)

	group := keeper.GetPosminingGroup(ctx.WithBlockTime(currentTime), paraminig, coin,balance)

	assert.Equal(t, len(group.Periods), 1)
	assert.Equal(t, group.Paramined, sdk.NewInt(0))
}

func TestKeeper_CalculatePosmined(t *testing.T) {
	// Make sure we won't get any posmined coins after threshold
	ctx, _, keeper, _, _ := createTestInput(t, false)

	coin := coins.Coin{
		Symbol: "test",
		PosminingThreshold: sdk.NewInt(1),
	}

	assert.Equal(t, keeper.CalculatePosmined(ctx, types.Posmining{}, coin, sdk.NewInt(1)), sdk.NewInt(0))

	coin = coins.Coin{
		Symbol: "test",
		PosminingThreshold: sdk.NewInt(0),
	}

	currentTime := MinusDays(0)
	LastCharge := MinusDays(1)
	LastTx := MinusDays(2)

	balance := sdk.NewIntWithDecimal(1000, 6)

	posmining := types.Posmining{
		DailyPercent:    sdk.NewInt(16), // 0.16% per day or 1.6 if balance 1000
		StructureCoff:   sdk.NewInt(218),
		Paramined:       sdk.NewInt(0),
		LastTransaction: LastTx,
		LastCharged:     LastCharge,
	}

	// 0.0016 per day * 2.18 structure coff
	shouldGetPosmined := balance.Mul(sdk.NewInt(16)).Quo(sdk.NewInt(10000)).
		Mul(sdk.NewInt(218)).Quo(sdk.NewInt(100))

	// 0.2
	group := keeper.GetPosminingGroup(ctx.WithBlockTime(currentTime), posmining, coin, balance)

	assert.Equal(t, len(group.Periods), 1)
	assert.Equal(t, group.Paramined, shouldGetPosmined)
}