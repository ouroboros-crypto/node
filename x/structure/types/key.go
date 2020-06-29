package types

const (
	// ModuleName is the name of the module
	ModuleName = "structure"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// Used for storing the upper-structure links
	FastAccessKey = "upper-structure"

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName
)
