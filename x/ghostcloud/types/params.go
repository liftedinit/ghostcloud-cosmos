package types

import (
	"gopkg.in/yaml.v2"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

const (
	DefaultMaxPayloadSize      int64  = 1024 * 1024 * 5 // 5MB
	DefaultMaxNameSize         int64  = 12
	DefaultMaxDescriptionSize  int64  = 512
	DefaultMaxUncompressedSize uint64 = 1024 * 1024 * 50 // 50MB
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		MaxPayloadSize:      DefaultMaxPayloadSize,
		MaxNameSize:         DefaultMaxNameSize,
		MaxDescriptionSize:  DefaultMaxDescriptionSize,
		MaxUncompressedSize: DefaultMaxUncompressedSize,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
