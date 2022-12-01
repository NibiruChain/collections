package examples

import (
	"github.com/NibiruChain/collections"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gogo/protobuf/proto"
)

type interfaceValueEncoder[T proto.Message] struct {
	cdc codec.BinaryCodec
}

func (a interfaceValueEncoder[T]) Encode(value T) []byte {
	accountAny, err := a.cdc.MarshalInterface(value)
	if err != nil {
		panic(err)
	}
	return accountAny
}

func (a interfaceValueEncoder[T]) Decode(b []byte) T {
	var acc T
	err := a.cdc.UnmarshalInterface(b, &acc)
	if err != nil {
		panic(err)
	}
	return acc
}

func (a interfaceValueEncoder[T]) Stringify(value T) string {
	return value.String()
}

func (a interfaceValueEncoder[T]) Name() string {
	return "auth.AccountInterface"
}

func NewInterfaceValueEncoder[T proto.Message](cdc codec.BinaryCodec) collections.ValueEncoder[T] {
	return interfaceValueEncoder[T]{cdc: cdc}
}
