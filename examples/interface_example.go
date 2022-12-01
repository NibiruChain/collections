package examples

import (
	"github.com/NibiruChain/collections"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type accountIValueEncoder struct {
	cdc codec.BinaryCodec
}

func (a accountIValueEncoder) Encode(value authtypes.AccountI) []byte {
	accountAny, err := a.cdc.MarshalInterface(value)
	if err != nil {
		panic(err)
	}
	return accountAny
}

func (a accountIValueEncoder) Decode(b []byte) authtypes.AccountI {
	var acc authtypes.AccountI
	err := a.cdc.UnmarshalInterface(b, &acc)
	if err != nil {
		panic(err)
	}
	return acc
}

func (a accountIValueEncoder) Stringify(value authtypes.AccountI) string {
	return value.String()
}

func (a accountIValueEncoder) Name() string {
	return "auth.AccountInterface"
}

func NewAccountInterfaceEncoder(cdc codec.BinaryCodec) collections.ValueEncoder[authtypes.AccountI] {
	return accountIValueEncoder{cdc: cdc}
}
