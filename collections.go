package collections

import (
	"context"
	"errors"

	"github.com/gogo/protobuf/io"
)

// ErrNotFound is returned when an object is not found.
var ErrNotFound = errors.New("collections: not found")

// Prefix defines a storage namespace which must be unique in a single module
// for all the different storage layer types: Map, Sequence, KeySet, Item, MultiIndex, IndexedMap
type Prefix uint8

func (n Prefix) Prefix() []byte { return []byte{uint8(n)} }

type Collection interface {
	Descriptor() CollectionDescriptor
	InitGenesis(context.Context, io.ReadCloser) error
	ExportGenesis(context.Context, io.WriteCloser) error
	Decode(key, value []byte) (k, v any, err error)
}

// KeyEncoder defines a generic interface which is implemented
// by types that are capable of encoding and decoding collections keys.
type KeyEncoder[T any] interface {
	// Encode encodes the type T into bytes.
	Encode(key T) []byte
	// Decode decodes the given bytes back into T.
	// And it also must return the bytes of the buffer which were read.
	Decode(b []byte) (int, T)
	// Stringify returns a string representation of T.
	Stringify(key T) string
	// Type describes the type of key encoder.
	Type() string
	//// Params describes parameters to the key encoder.
	//Params() []interface{}
}

// ValueEncoder defines a generic interface which is implemented
// by types that are capable of encoding and decoding collection values.
type ValueEncoder[T any] interface {
	// Encode encodes the value T into bytes.
	Encode(value T) []byte
	// Decode returns the type T given its bytes representation.
	Decode(b []byte) T
	// Stringify returns a string representation of T.
	Stringify(value T) string
	// Type returns the type of the object.
	Type() string
}
