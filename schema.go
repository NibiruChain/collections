package collections

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Schema struct {
	storeKey   sdk.StoreKey
	descriptor *SchemaDescriptor
	namespaces map[Namespace]bool
	names      map[string]bool
}

type SchemaDescriptor struct {
	Maps      []MapDescriptor      `json:"maps,omitempty"`
	Sequences []SequenceDescriptor `json:"sequences,omitempty"`
	Items     []ItemDescriptor     `json:"items,omitempty"`
}

type MapDescriptor struct {
	Prefix []byte          `json:"prefix"`
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

type SequenceDescriptor struct {
	Prefix []byte `json:"prefix"`
	Name   string `json:"name"`
}

type ItemDescriptor struct {
	Prefix []byte
	Name   string
	Type   string
}

func NewSchema(storeKey sdk.StoreKey) Schema {
	return Schema{
		storeKey:   storeKey,
		descriptor: &SchemaDescriptor{},
		namespaces: map[Namespace]bool{},
		names:      map[string]bool{},
	}
}

func (s Schema) Descriptor() SchemaDescriptor {
	return *s.descriptor
}

func (s Schema) ensureUniqueNamespace(namespace Namespace) {
	if s.namespaces[namespace] {
		panic(fmt.Errorf("namespace %d already taken within schema", namespace))
	}

	s.namespaces[namespace] = true
}

func (s Schema) ensureUniqueName(name string) {
	// TODO valid name format
	if s.names[name] {
		panic(fmt.Errorf("name %d already taken within schema", name))
	}

	s.names[name] = true
}
