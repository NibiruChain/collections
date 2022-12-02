package examples

import (
	"github.com/cosmos/cosmos-sdk/codec"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/NibiruChain/collections"
)

type AccountKeeper struct {
	Schema        collections.Schema
	AccountNumber collections.Sequence
	Accounts      collections.Map[sdk.AccAddress, types.BaseAccount]
	Params        collections.Item[types.Params]
}

func NewAccountKeeper(sk sdk.StoreKey, cdc codec.BinaryCodec) *AccountKeeper {
	schema := collections.NewSchema(sk)
	return &AccountKeeper{
		Schema:        schema,
		AccountNumber: collections.NewSequence(schema, 0, "account_number_seq"),
		Accounts: collections.NewMap(schema, 1,
			"address", collections.AccAddressKeyEncoder,
			"account", collections.ProtoValueEncoder[types.BaseAccount](cdc),
		),
		Params: collections.NewItem(schema, 2, "params", collections.ProtoValueEncoder[types.Params](cdc)),
	}
}

func (k AccountKeeper) CreateAccount(ctx sdk.Context, pubKey crypto.PubKey) {
	n := k.AccountNumber.Next(ctx)
	addr := sdk.AccAddress(pubKey.String())
	acc := types.BaseAccount{
		Address:       addr.String(),
		AccountNumber: n,
		Sequence:      0,
	}

	k.Accounts.Insert(ctx, addr, acc)
}

func (k AccountKeeper) GetAccount(ctx sdk.Context, addr sdk.AccAddress) (types.BaseAccount, error) {
	return k.Accounts.Get(ctx, addr)
}

func (k AccountKeeper) AllAccounts(ctx sdk.Context) []types.BaseAccount {
	return k.Accounts.Iterate(ctx, collections.Range[sdk.AccAddress]{}).Values()
}

func (k AccountKeeper) AllAddresses(ctx sdk.Context) []sdk.AccAddress {
	return k.Accounts.Iterate(ctx, collections.Range[sdk.AccAddress]{}).Keys()
}

func (k AccountKeeper) AccountsBetween(ctx sdk.Context, start, end sdk.AccAddress) []types.BaseAccount {
	rng := collections.Range[sdk.AccAddress]{}.
		StartInclusive(start).
		EndInclusive(end)
	// .Descending() for reverse order
	return k.Accounts.Iterate(ctx, rng).Values()
}
