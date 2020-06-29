package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Мы используем эту структуру как ссылку на владельца верхней структуры
type UpperStructure struct {
	Owner   sdk.AccAddress `json:"owner"`
	Address sdk.AccAddress `json:"address"`
}

// Структура пользователя
type Structure struct {
	Owner     sdk.AccAddress `json:"owner"`
	Balance   sdk.Int          `json:"balance"` // баланс структуры
	Followers sdk.Int `json:"followers"` // кол-во последователей
	MaxLevel  sdk.Int `json:"max_level"` // максимальный уровень структуры
}

func (r UpperStructure) String() string {
	return r.Owner.String()
}

func (r Structure) String() string {
	return r.Balance.String()
}

// Возвращает новую UpperStructure
func NewUpperStructure(address sdk.AccAddress) UpperStructure {
	return UpperStructure{
		Address: address,
	}
}

func NewStructure(owner sdk.AccAddress) Structure {
	return Structure{
		Owner: owner,
		Balance: sdk.NewInt(0),
		Followers: sdk.NewInt(0),
		MaxLevel: sdk.NewInt(0),
	}
}

