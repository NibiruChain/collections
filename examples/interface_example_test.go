package examples

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	db "github.com/tendermint/tm-db"
	"testing"
)

func deps() (sdk.StoreKey, sdk.Context, codec.BinaryCodec) {
	sk := sdk.NewKVStoreKey("mock")
	dbm := db.NewMemDB()
	ms := store.NewCommitMultiStore(dbm)
	ms.MountStoreWithDB(sk, types.StoreTypeIAVL, dbm)
	if err := ms.LoadLatestVersion(); err != nil {
		panic(err)
	}

	ir := codectypes.NewInterfaceRegistry()
	authtypes.RegisterInterfaces(ir)
	return sk,
		sdk.Context{}.
			WithMultiStore(ms).
			WithGasMeter(sdk.NewGasMeter(1_000_000_000)),
		codec.NewProtoCodec(ir)
}

func Test_something(t *testing.T) {
	sk, ctx, cdc := deps()
	k := NewInterfaceKeeper(sk, cdc)

	// using interface type
	var genericAccount authtypes.AccountI = &authtypes.BaseAccount{}
	k.Accounts.Insert(ctx, sdk.AccAddress("generic account"), genericAccount)

	// using concrete types
	k.Accounts.Insert(ctx, sdk.AccAddress("module"), &authtypes.ModuleAccount{
		BaseAccount: &authtypes.BaseAccount{}, // NOTE: required otherwise it will panic because of bug in ModuleAccount.UnpackAny
		Name:        "module account",
	})

	k.Accounts.Insert(ctx, sdk.AccAddress("external user"), &authtypes.BaseAccount{
		PubKey:        nil,
		AccountNumber: 1,
		Sequence:      2,
	})

	// getting interface type
	moduleAcc, err := k.Accounts.Get(ctx, sdk.AccAddress("module"))
	require.NoError(t, err)

	_ = moduleAcc.(*authtypes.ModuleAccount)

	userAcc, err := k.Accounts.Get(ctx, sdk.AccAddress("external user"))
	require.NoError(t, err)

	_ = userAcc.(*authtypes.BaseAccount)
}