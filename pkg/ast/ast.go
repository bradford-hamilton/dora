// Package ast TODO: package docs
package ast

// These are the available root node types. In JSON it will either be an
// object or an array at the base.
const (
	ObjectRoot RootNodeType = iota
	ArrayRoot
)

type RootNodeType int

// RootNode is what starts every parsed AST. There is a `Type` field so that
// you can ask which root node type starts the tree. From there you can access the
type RootNode struct {
	RootValue *Value
	Type      RootNodeType
}

// Object represents a JSON object. It holds a slice of Property as its children.
type Object struct {
	Type     string // "Object"
	Children []Property
}

// Literal represents a JSON literal value. It holds it's Type as well as value.
type Literal struct {
	Type  string // "Literal"
	Value Value
}

// Array represents a JSON array value. It holds it's Type as well as a slice of children values.
type Array struct {
	Type     string // "Array"
	Children []Value
}

// Property holds its own Type as well as a `Key` and `Value`. The Key is an Identifier
// and the value is a Value so that we can continue to parse obj/array/literals at this point.
type Property struct {
	Type  string
	Key   Identifier
	Value Value
}

// Identifier represents a JSON object property key
type Identifier struct {
	Type  string // "Identifier"
	Value string // "key1"
	Raw   string // "\"key1\""
}

// Value will eventually hold some methods that all Values must implement. For now
// it is what allows us to switch over 3 different parsable types
type Value interface{}

// Available object states
const (
	ObjStart objectState = iota
	ObjOpen
	ObjProperty
	ObjComma
)

type objectState int

// Available property states
const (
	PropertyStart propertyState = iota
	PropertyKey
	PropertyColon
)

type propertyState int

// Available array states
const (
	ArrayStart arrayState = iota
	ArrayOpen
	ArrayValue
	ArrayComma
)

type arrayState int

// Available string states
const (
	StringStart stringState = iota
	StringQuoteOrChar
	Escape
)

type stringState int

// Available number states
const (
	NumberStart numberState = iota
	NumberMinus
	NumberZero
	NumberDigit
	NumberPoint
	NumberDigitFraction
	NumberExp
	NumberExpDigitOrSign
)

type numberState int
