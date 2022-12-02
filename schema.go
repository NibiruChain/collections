package collections

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Schema struct {
	storeKey            sdk.StoreKey
	descriptor          *SchemaDescriptor
	collectionsByPrefix map[string]Collection
	collectionsByName   map[string]Collection
}

type SchemaDescriptor struct {
	Collections []CollectionDescriptor `json:"collections,omitempty"`
}

type CollectionDescriptor struct {
	Type   string          `json:"type"`
	Prefix []byte          `json:"prefix"`
	Name   string          `json:"name"`
	Key    KeyDescriptor   `json:"key"`
	Value  ValueDescriptor `json:"value"`
}

type KeyDescriptor struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ValueDescriptor struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func NewSchema(storeKey sdk.StoreKey) Schema {
	return Schema{
		storeKey:            storeKey,
		descriptor:          &SchemaDescriptor{},
		collectionsByPrefix: map[string]Collection{},
		collectionsByName:   map[string]Collection{},
	}
}

func (s Schema) Descriptor() SchemaDescriptor {
	return *s.descriptor
}

func (s Schema) addCollection(collection Collection) {
	desc := collection.Descriptor()
	prefix := desc.Prefix
	name := desc.Name

	if _, ok := s.collectionsByPrefix[string(prefix)]; ok {
		panic(fmt.Errorf("prefix %v already taken within schema", desc.Prefix))
	}

	if _, ok := s.collectionsByName[name]; ok {
		panic(fmt.Errorf("prefix %d already taken within schema", name))
	}

	if !nameRegex.MatchString(name) {
		panic(fmt.Errorf("name must match regex %s, got %s", nameRegex.String(), name))
	}

	s.collectionsByPrefix[string(prefix)] = collection
	s.collectionsByName[name] = collection
	s.descriptor.Collections = append(s.descriptor.Collections, desc)
}

var nameRegex = regexp.MustCompile("^[A-Za-z][A-Za-z0-9_]*$")

func (s Schema) InitGenesis(ctx context.Context, message json.RawMessage) error {
	panic("TODO")
}

func (s Schema) ExportGenesis(ctx context.Context) (json.RawMessage, error) {
	panic("TODO")
}

// Decode implements logical decoding.
func (s Schema) Decode(key, value []byte) (Entry, error) {
	panic("TODO")
}

// Entry is a logical decoding of key value pair bytes.
type Entry struct {
	CollectionName string
	Key            any
	Value          any
}
