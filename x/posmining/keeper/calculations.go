package keeper

import (
	"github.com/ouroboros-crypto/node/x/posmining/types"
	"math"
	"time"
)

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ouroboros-crypto/node/x/coins"
)

// Returns a list of saving posmining periods
func (k Keeper) GetSavingPeriods(ctx sdk.Context, posmining types.Posmining) []types.PosminingPeriod {
	var daysSeparator int64 = 2592000

	lastTx := posmining.LastTransaction
	secondsDiff := int64(ctx.BlockTime().Sub(lastTx).Seconds())

	if secondsDiff < daysSeparator {
		return []types.PosminingPeriod{types.NewPosminingPeriod(lastTx, ctx.BlockTime(), sdk.NewInt(0), sdk.NewInt(0))}
	}

	updateThreshold, _ := time.Parse(time.RFC822, "09 Aug 21 12:00 UTC")

	// Saving coff won't be working after 09 Aug 21
	if ctx.BlockHeader().Time.After(updateThreshold) {
		return []types.PosminingPeriod{types.NewPosminingPeriod(lastTx, ctx.BlockTime(), sdk.NewInt(0), sdk.NewInt(0))}
	}

	periods := secondsDiff / daysSeparator
	mod := int64(math.Mod(float64(secondsDiff), float64(daysSeparator)))

	var result []types.PosminingPeriod
	var i int64 = 0

	for i < periods {
		result = append(result, types.NewPosminingPeriod(
			lastTx.Add(time.Duration(daysSeparator*i)*time.Second),
			lastTx.Add(time.Duration(daysSeparator*(i+1))*time.Second),
			sdk.NewInt(0),
			types.GetSavingCoff(int(i)),
		))

		i += 1
	}

	// What's left
	if mod > 0 {
		latestPeriod := lastTx.Add(time.Duration(daysSeparator*periods) * time.Second)

		result = append(result, types.NewPosminingPeriod(
			latestPeriod,
			latestPeriod.Add(time.Duration(mod)*time.Second),
			sdk.NewInt(0),
			types.GetSavingCoff(int(periods)),
		))
	}

	return result
}

// Returns a list of correction posmining periods
func (k Keeper) GetCorrectionPeriods(ctx sdk.Context, posmining types.Posmining) []types.PosminingPeriod {
	correction := k.GetCorrection(ctx)

	// First we always initialize the current correction period
	result := []types.PosminingPeriod{
		types.NewPosminingPeriod(correction.StartDate, ctx.BlockTime(), correction.CorrectionCoff, sdk.NewInt(0)),
	}

	// If we should count only the current period
	if correction.StartDate.Before(posmining.LastCharged) {
		result[0].Start = posmining.LastCharged

		return result
	}

	for _, previous := range correction.PreviousCorrections {
		if previous.EndDate.After(posmining.LastCharged) {
			var startDate time.Time

			if previous.StartDate.Before(posmining.LastCharged) {
				startDate = posmining.LastCharged
			} else {
				startDate = previous.StartDate
			}

			result = append([]types.PosminingPeriod{types.NewPosminingPeriod(startDate, previous.EndDate, previous.CorrectionCoff, sdk.NewInt(0))}, result...)
		}

		if previous.StartDate.Before(posmining.LastCharged) {
			break
		}
	}

	return result
}

// Calculates and returns a group of posmining periods
func (k Keeper) GetPosminingGroup(ctx sdk.Context, posmining types.Posmining, coin coins.Coin, balance sdk.Int) types.PosminingGroup {
	group := types.NewPosminingGroup(posmining, balance)

	// For the custom coins, we just have to apply the usual percents during the whole time
	if !coin.Default {
		group.Add(types.NewPosminingPeriod(posmining.LastCharged, ctx.BlockTime(), sdk.NewInt(0), sdk.NewInt(0)))

		return group
	}

	savings := k.GetSavingPeriods(ctx, posmining)

	corrections := k.GetCorrectionPeriods(ctx, posmining)

	r_i := 0

	for _, saving := range savings {
		for r_i < len(corrections) {
			correction := corrections[r_i]

			// In that case, we should move to the next saving
			if correction.Start.Equal(saving.End) || correction.Start.After(saving.End) {
				break
			}

			if correction.End.Equal(saving.End) || correction.End.Before(saving.End) {
				correction.SavingCoff = saving.SavingCoff

				group.Add(correction)
			} else {
				newCorrection := correction

				correction.End = saving.End
				correction.SavingCoff = saving.SavingCoff

				group.Add(correction)

				newCorrection.Start = correction.End

				// Making new space
				corrections = append(corrections, types.PosminingPeriod{})

				// Copying element to the next index
				copy(corrections[r_i+2:], corrections[r_i+1:])

				// Replacing
				corrections[r_i+1] = newCorrection
			}

			r_i += 1
		}
	}

	return group
}

// Calculates how many tokens has been posmined
func (k Keeper) CalculatePosmined(ctx sdk.Context, posmining types.Posmining, coin coins.Coin, coinsAmount sdk.Int) sdk.Int {
	updateThreshold, _ := time.Parse(time.RFC822, "09 Aug 21 12:00 UTC")

	// If we have a threshold set and it's already has been reach, we should always return 0
	if ctx.BlockHeader().Time.Before(updateThreshold) && coin.PosminingThreshold.IsPositive() && coinsAmount.GTE(coin.PosminingThreshold) {
		return sdk.NewInt(0)
	}

	// If posmining was disabled by the owner
	if ctx.BlockHeader().Time.After(updateThreshold) && !k.GetPosminingEnabled(ctx, posmining.Owner) {
		return sdk.NewInt(0)
	}

	posmined := posmining.Paramined.Add(k.GetPosminingGroup(ctx, posmining, coin, coinsAmount).Paramined)

	if coin.PosminingThreshold.IsPositive() && posmined.IsPositive() && coinsAmount.Add(posmined).GTE(coin.PosminingThreshold) {
		posmined = coin.PosminingThreshold.Sub(coinsAmount)
	}

	return posmined
}
