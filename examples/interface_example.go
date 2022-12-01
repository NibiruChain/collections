package examples

import (
	"github.com/NibiruChain/collections"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

type MyInterfaceKeeper struct {
	Accounts collections.Map[sdk.AccAddress, authtypes.AccountI]
}

func NewInterfaceKeeper(sk sdk.StoreKey, cdc codec.BinaryCodec) MyInterfaceKeeper {
	return MyInterfaceKeeper{
		Accounts: collections.NewMap(sk, 0, collections.AccAddressKeyEncoder, NewAccountInterfaceEncoder(cdc)),
	}
}
