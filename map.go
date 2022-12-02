package collections

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Map[K, V any] struct {
	kc KeyEncoder[K]
	vc ValueEncoder[V]

	prefix []byte
	sk     sdk.StoreKey

	keyName  string
	name     string
	typeName string
}

func NewMap[K, V any](schema Schema, namespace Prefix,
	keyName string, kc KeyEncoder[K],
	valueName string, vc ValueEncoder[V]) Map[K, V] {
	m := newMap(schema.storeKey, namespace, keyName, kc, valueName, vc)
	schema.addCollection(m)
	return m
}

func newMap[K, V any](sk sdk.StoreKey, namespace Prefix,
	keyName string,
	kc KeyEncoder[K],
	valueName string,
	vc ValueEncoder[V]) Map[K, V] {
	return Map[K, V]{
		kc:     kc,
		vc:     vc,
		prefix: namespace.Prefix(),
		sk:     sk,
		//nolint
		typeName: vc.(ValueEncoder[V]).Type(), // go1.19 compiler bug
		name:     valueName,
		keyName:  keyName,
	}
}

func (m Map[K, V]) Insert(ctx sdk.Context, k K, v V) {
	m.getStore(ctx).
		Set(m.kc.Encode(k), m.vc.Encode(v))
}

func (m Map[K, V]) Get(ctx sdk.Context, k K) (v V, err error) {
	vBytes := m.getStore(ctx).Get(m.kc.Encode(k))
	if vBytes == nil {
		return v, fmt.Errorf("%w: '%s' with key %s", ErrNotFound, m.typeName, m.kc.Stringify(k))
	}

	return m.vc.Decode(vBytes), nil
}

func (m Map[K, V]) GetOr(ctx sdk.Context, key K, def V) (v V) {
	v, err := m.Get(ctx, key)
	if err == nil {
		return
	}

	return def
}

func (m Map[K, V]) Delete(ctx sdk.Context, k K) error {
	kBytes := m.kc.Encode(k)
	store := m.getStore(ctx)
	if !store.Has(kBytes) {
		return fmt.Errorf("%w: '%s' with key %s", ErrNotFound, m.typeName, m.kc.Stringify(k))
	}
	store.Delete(kBytes)

	return nil
}

func (m Map[K, V]) Iterate(ctx sdk.Context, rng Ranger[K]) Iterator[K, V] {
	return iteratorFromRange[K, V](m.getStore(ctx), rng, m.kc, m.vc)
}

func (m Map[K, V]) getStore(ctx sdk.Context) sdk.KVStore {
	return prefix.NewStore(ctx.KVStore(m.sk), m.prefix)
}

func (m Map[K, V]) Descriptor() CollectionDescriptor {
	return CollectionDescriptor{
		Type:   "map",
		Prefix: m.prefix,
		Name:   m.name,
		Key: KeyDescriptor{
			Name: m.keyName,
			Type: m.kc.Type(),
		},
		Value: ValueDescriptor{
			Name: m.name,
			Type: m.vc.Type(),
		},
	}
}
