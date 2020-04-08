// Package ast TODO: package docs
package ast

// TODO: docs
const (
	ObjectRoot rootNodeType = iota
	ArrayRoot
)

type rootNodeType int

// RootNode TODO
type RootNode struct {
	Object *Object
	Array  *Array
	Type   rootNodeType
}

// Object TODO
type Object struct {
	Type     string // "Object"
	Children []Property
}

// Literal TODO
type Literal struct {
	Type  string // "Literal"
	Value Value
}

// Array TODO
type Array struct {
	Type     string // "Array"
	Children []Value
}

// ============================================

// Property TODO
type Property struct {
	Type  string
	Key   Identifier
	Value Value
}

// Identifier TODO
type Identifier struct {
	Type  string // "Identifier"
	Value string // "key1"
	Raw   string // "\"key1\""
}

// Value TODO
type Value interface{}
