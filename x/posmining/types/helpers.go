package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math"
	"time"
)

// Для подсчета начисляемых токенов за какое-то время (сутки, час, минута, секунда)
type CoinsPerTime struct {
	Day    sdk.Int `json:"day"`
	Hour   sdk.Int `json:"hour"`
	Minute sdk.Int `json:"minute"`
	Second sdk.Int `json:"second"`
}

// Calculates and returns new CoinsPerTime
func NewCoinsPerTime(balance sdk.Int, dailyPercent sdk.Int, structureCoff sdk.Int, savingsCoff sdk.Int, regulationCOff sdk.Int) CoinsPerTime {
	result := CoinsPerTime{
		Day:    sdk.NewInt(0),
		Hour:   sdk.NewInt(0),
		Minute: sdk.NewInt(0),
		Second: sdk.NewInt(0),
	}

	toQuo := sdk.NewInt(10000)

	actualPercent := dailyPercent

	if structureCoff.IsZero() == false {
		actualPercent = actualPercent.Mul(structureCoff)
		toQuo = toQuo.MulRaw(100)
	}

	if savingsCoff.IsZero() == false {
		actualPercent = actualPercent.Mul(savingsCoff)
		toQuo = toQuo.MulRaw(100)
	}

	if regulationCOff.IsZero() == false {
		actualPercent = actualPercent.Mul(regulationCOff)
		toQuo = toQuo.MulRaw(100)
	}

	result.Day = balance.Mul(actualPercent).Quo(toQuo)
	result.Hour = result.Day.QuoRaw(24)
	result.Minute = result.Hour.QuoRaw(60)
	result.Second = result.Minute.QuoRaw(60)

	return result
}


// Time difference
type TimeDifference struct {
	Days    sdk.Int `json:"days"`    // Кол-во days
	Hours   sdk.Int `json:"hours"`   // Кол-во часов
	Minutes sdk.Int `json:"minutes"` // Кол-во минут
	Seconds sdk.Int `json:"seconds"` // Кол-во секунд

	Total sdk.Int `json:"total"` // Общее время в секундах
}

// Creates new time difference based on the seconds difference
func NewTimeDifference(seconds sdk.Int) TimeDifference {
	difference := TimeDifference{
		Days:    sdk.NewInt(0),
		Hours:   sdk.NewInt(0),
		Minutes: sdk.NewInt(0),
		Seconds: sdk.NewInt(0),
		Total:   sdk.NewInt(0),
	}

	duration := time.Duration(seconds.Int64()) * time.Second

	// Less then a minute
	if duration.Seconds() < 60.0 {
		difference.Seconds = sdk.NewInt(int64(duration.Seconds()))

		return difference
	}

	// Less than an hour
	if duration.Minutes() < 60.0 {
		difference.Minutes = sdk.NewInt(int64(duration.Minutes()))
		difference.Seconds = sdk.NewInt(int64(math.Mod(duration.Seconds(), 60)))

		return difference
	}

	// Less than a day
	if duration.Hours() < 24.0 {
		difference.Hours = sdk.NewInt(int64(duration.Hours()))
		difference.Minutes = sdk.NewInt(int64(math.Mod(duration.Minutes(), 60)))
		difference.Seconds = sdk.NewInt(int64(math.Mod(duration.Seconds(), 60)))

		return difference
	}

	difference.Days = sdk.NewInt(int64(duration.Hours() / 24))
	difference.Hours = sdk.NewInt(int64(math.Mod(duration.Hours(), 24)))
	difference.Minutes = sdk.NewInt(int64(math.Mod(duration.Minutes(), 60)))
	difference.Seconds = sdk.NewInt(int64(math.Mod(duration.Seconds(), 60)))

	return difference
}


// A single posmining period
type PosminingPeriod struct {
	Start          time.Time `json:"start"`       // Начало периода
	End            time.Time `json:"end"`         // конец периода
	CorrectionCoff sdk.Int   `json:"regulation"`  // Регуляция
	SavingCoff     sdk.Int   `json:"saving_coff"` // Коэффициент накопления
}

// Hoe much time pass between Start and End
func (p PosminingPeriod) TimePass() TimeDifference {
	return NewTimeDifference(sdk.NewInt(int64(p.End.Sub(p.Start).Seconds())))
}

func NewPosminingPeriod(start time.Time, end time.Time, regulationCoff sdk.Int, savingCoff sdk.Int) PosminingPeriod {
	return PosminingPeriod{
		start,
		end,
		regulationCoff,
		savingCoff,
	}
}

// A group of posmining periods
type PosminingGroup struct {
	Paramined  sdk.Int           `json:"paramined"`  // How many coins paramined during period
	Balance    sdk.Int           `json:"balance"`    // Current balance
	Posmining Posmining        `json:"posmining"` // Current balance
	Periods    []PosminingPeriod `json:"periods"`    // Posmining periods
}

func NewPosminingGroup(posmining Posmining, balance sdk.Int) PosminingGroup {
	return PosminingGroup{
		Paramined: sdk.NewInt(0),
		Balance: balance,
		Posmining: posmining,
	}
}

// Adds a posmining period
func (p *PosminingGroup) Add(period PosminingPeriod) {
	perTime := NewCoinsPerTime(p.Balance, p.Posmining.DailyPercent, p.Posmining.StructureCoff, period.SavingCoff, period.CorrectionCoff)

	timeDiff := period.TimePass()

	// Add to the total
	p.Paramined = p.Paramined.Add(sdk.NewInt(0).Add(timeDiff.Seconds.Mul(
		perTime.Second).Add(
		timeDiff.Minutes.Mul(perTime.Minute)).Add(
		timeDiff.Hours.Mul(perTime.Hour)).Add(
		timeDiff.Days.Mul(perTime.Day))))

	// Append to the periods
	p.Periods = append(p.Periods, period)

}
