package bank

import (
	"github.com/ouroboros-crypto/node/x/bank/keeper"
)

var (
	NewKeeper = keeper.NewKeeper
)

type (
	Keeper = keeper.Keeper
)
