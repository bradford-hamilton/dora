// Package ast TODO: package docs
package ast

// These are the available root node types. In JSON it will either be an
// object or an array at the base.
const (
	ObjectRoot RootNodeType = iota
	ArrayRoot
)

// RootNodeType is a type alias for an int
type RootNodeType int

// RootNode is what starts every parsed AST. There is a `Type` field so that
// you can ask which root node type starts the tree.
type RootNode struct {
	RootValue *Value
	Type      RootNodeType
}

// Object represents a JSON object. It holds a slice of Property as its children,
// a Type ("Object"), and start & end code points for displaying.
type Object struct {
	Type     string // "Object"
	Children []Property
	Start    int
	End      int
}

// Array represents a JSON array It holds a slice of Value as its children,
// a Type ("Array"), and start & end code points for displaying.
type Array struct {
	Type     string // "Array"
	Children []Value
	Start    int
	End      int
}

// Literal represents a JSON literal value. It holds a Type ("Literal") and the actual value.
type Literal struct {
	Type  string // "Literal"
	Value Value
}

// Property holds a Type ("Property") as well as a `Key` and `Value`. The Key is an Identifier
// and the value is any Value.
type Property struct {
	Type  string // "Property"
	Key   Identifier
	Value Value
}

// Identifier represents a JSON object property key
type Identifier struct {
	Type  string // "Identifier"
	Value string // "key1"
}

// Value will eventually have some methods that all Values must implement. For now
// it represents any JSON value (object | array | boolean | string | number | null)
type Value interface{}

// state is a type alias for int and used to create the available value states below
type state int

// Available states for each type used in parsing
const (
	// Object states
	ObjStart state = iota
	ObjOpen
	ObjProperty
	ObjComma

	// Property states
	PropertyStart
	PropertyKey
	PropertyColon

	// Array states
	ArrayStart
	ArrayOpen
	ArrayValue
	ArrayComma

	// String states
	StringStart
	StringQuoteOrChar
	Escape

	// Number states
	NumberStart
	NumberMinus
	NumberZero
	NumberDigit
	NumberPoint
	NumberDigitFraction
	NumberExp
	NumberExpDigitOrSign
)
