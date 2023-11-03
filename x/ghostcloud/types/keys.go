package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName defines the module name
	ModuleName = "ghostcloud"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_ghostcloud"
)

var (
	DeploymentMetaKeyPrefix     = []byte{0x00}
	DeploymentItemKeyPrefix     = []byte{0x01}
	DeploymentItemMetaPrefix    = []byte{DeploymentItemKeyPrefix[0], 0x00}
	DeploymentItemContentPrefix = []byte{DeploymentItemKeyPrefix[0], 0x01}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func DeploymentItemKey(
	addr sdk.AccAddress,
	name string,
	file string,
) []byte {
	var key []byte

	addrBytes := []byte(addr)
	nameBytes := []byte(name)
	key = append(key, addrBytes...)
	key = append(key, nameBytes...)
	key = append(key, file...)

	return key
}

func DeploymentKey(
	addr sdk.AccAddress,
	name string,
) []byte {
	var key []byte

	addrBytes := []byte(addr)
	nameBytes := []byte(name)
	key = append(key, addrBytes...)
	key = append(key, nameBytes...)

	return key
}
