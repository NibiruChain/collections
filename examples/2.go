package examples

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/NibiruChain/collections"
)

// let's showcase some more complex keys, like delegations which are composite.

type StakingKeeper struct {
	Schema      collections.Schema
	Delegations collections.Map[collections.Pair[sdk.ValAddress, sdk.AccAddress], types.Delegation]
}

func NewStakingKeeper(sk sdk.StoreKey, cdc codec.BinaryCodec) *StakingKeeper {
	schema := collections.NewSchema(sk)
	return &StakingKeeper{
		Schema: schema,
		Delegations: collections.NewMap(
			schema, 0,
			"val_address_acc_address", collections.PairKeyEncoder(collections.ValAddressKeyEncoder, collections.AccAddressKeyEncoder), // we pass here a joint key encoder which encodes both val address key and acc address key
			"delegation", collections.ProtoValueEncoder[types.Delegation](cdc),
		),
	}
}

func (k StakingKeeper) CreateDelegation(ctx sdk.Context, val sdk.ValAddress, del sdk.AccAddress) {
	k.Delegations.Insert(ctx, collections.Join(val, del), types.Delegation{
		DelegatorAddress: del.String(),
		ValidatorAddress: val.String(),
		Shares:           sdk.MustNewDecFromStr("100000"),
	})
}

func (k StakingKeeper) GetDelegation(ctx sdk.Context, val sdk.ValAddress, del sdk.AccAddress) (types.Delegation, error) {
	return k.Delegations.Get(ctx, collections.Join(val, del))
}

func (k StakingKeeper) GetValidatorDelegations(ctx sdk.Context, val sdk.ValAddress) []types.Delegation {
	rng := collections.PairRange[sdk.ValAddress, sdk.AccAddress]{}.
		Prefix(val) // gets all the keys starting with val [it's prefix safe]

	return k.Delegations.Iterate(ctx, rng).Values()
}

func (k StakingKeeper) GetValidatorDelegationsBetween(ctx sdk.Context, val sdk.ValAddress, start sdk.AccAddress, end sdk.AccAddress) []types.Delegation {
	rng := collections.PairRange[sdk.ValAddress, sdk.AccAddress]{}.
		Prefix(val). // gets all the keys starting with val [it's prefix safe]
		StartInclusive(start).
		EndInclusive(end)
	return k.Delegations.Iterate(ctx, rng).Values()
}
