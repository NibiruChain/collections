package collections

import (
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	"cosmossdk.io/store/types"
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func deps() (types.StoreKey, sdk.Context, codec.BinaryCodec) {
	sk := storetypes.NewKVStoreKey("mock")
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	ms.MountStoreWithDB(sk, types.StoreTypeIAVL, db)
	if err := ms.LoadLatestVersion(); err != nil {
		panic(err)
	}

	return sk,
		sdk.Context{}.
			WithMultiStore(ms).
			WithGasMeter(storetypes.NewGasMeter(1_000_000_000)),
		codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
}
