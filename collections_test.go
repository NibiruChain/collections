package collections

import (
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func deps() (storetypes.StoreKey, sdk.Context, codec.BinaryCodec) {
	sk := storetypes.NewKVStoreKey("mock")
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	ms.MountStoreWithDB(sk, storetypes.StoreTypeIAVL, db)
	if err := ms.LoadLatestVersion(); err != nil {
		panic(err)
	}

	return sk,
		sdk.Context{}.
			WithMultiStore(ms).
			WithGasMeter(storetypes.NewGasMeter(1_000_000_000)),
		codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
}
