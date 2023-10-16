package types

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

var _ binary.ByteOrder

var (
	// DeploymentKeyPrefix is the prefix to retrieve all Deployment
	DeploymentKeyPrefix = []byte{0x00}
)

// DeploymentKey returns the store key to retrieve a Deployment from the index fields
func DeploymentKey(
	addr sdk.AccAddress,
	name string,
) []byte {
	var key []byte

	addrBytes := []byte(addr)
	nameBytes := []byte(name)
	key = append(key, addrBytes...)
	key = append(key, nameBytes...)
	key = append(key, []byte("/")...)

	return key
}

func CreateAccountStorePrefix(addr sdk.AccAddress) []byte {
	return append(DeploymentKeyPrefix, address.MustLengthPrefix(addr)...)
}

func CreateAccountDeploymentPrefix(addr sdk.AccAddress, name string) []byte {
	return append(CreateAccountStorePrefix(addr), KeyPrefix(name)...)
}
