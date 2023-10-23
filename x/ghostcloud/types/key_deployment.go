package types

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ binary.ByteOrder

var (
	DeploymentMetaKeyPrefix      = []byte{0x00}
	DeploymentFilesKeyPrefix     = []byte{0x01}
	DeploymentFilesMetaPrefix    = []byte{DeploymentFilesKeyPrefix[0], 0x00}
	DeploymentFilesContentPrefix = []byte{DeploymentFilesKeyPrefix[0], 0x01}
)

func DeploymentFileKey(
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
