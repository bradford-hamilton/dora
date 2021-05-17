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

// Available ast value types
const (
	ObjectType Type = iota
	ArrayType
	ArrayItemType
	LiteralType
	PropertyType
	IdentifierType
)

const (
	StringLiteralValueType LiteralValueType = iota
	NumberLiteralValueType
	NullLiteralValueType
	BooleanLiteralValueType
)

// Type is a type alias for int. Represents a node's type.
type Type int

// LiteralValueType is a type alias for int. Represents the type of the value in a Literal node
type LiteralValueType int

type StructuralItem struct {
	Value string
}

// Object represents a JSON object. It holds a slice of Property as its children,
// a Type ("Object"), and start & end code points for displaying.
type Object struct {
	Type            Type
	Children        []Property
	Start           int
	End             int
	SuffixStructure []StructuralItem
}

// Array represents a JSON array It holds a slice of Value as its children,
// a Type ("Array"), and start & end code points for displaying.
type Array struct {
	Type            Type
	PrefixStructure []StructuralItem
	Children        []ArrayItem
	SuffixStructure []StructuralItem
	Start           int
	End             int
}

// Array holds a Type ("ArrayItem") as well as a `Value` and whether there is a comma after the item
type ArrayItem struct {
	Type               Type
	PrefixStructure    []StructuralItem
	Value              ValueContent
	PostValueStructure []StructuralItem
	HasCommaSeparator  bool
}

// Literal represents a JSON literal value. It holds a Type ("Literal") and the actual value.
type Literal struct {
	Type              Type
	ValueType         LiteralValueType
	Value             ValueContent
	Delimiter         string // Delimiter is set for string values
	OriginalRendering string // Allows preservig numeric formatting from source documents
}

// Property holds a Type ("Property") as well as a `Key` and `Value`. The Key is an Identifier
// and the value is any Value.
type Property struct {
	Type               Type
	PrefixStructure    []StructuralItem
	Key                Identifier
	PostKeyStructure   []StructuralItem // NOTE: Colon is between PostKeyStructure and PreValue Structure
	PreValueStructure  []StructuralItem
	Value              ValueContent
	PostValueStructure []StructuralItem
	HasCommaSeparator  bool
}

// Identifier represents a JSON object property key
type Identifier struct {
	Type      Type
	Value     string // "key1"
	Delimiter string
}

type Value struct {
	PrefixStructure []StructuralItem
	Content         ValueContent
	SuffixStructure []StructuralItem
}

// ValueContent will eventually have some methods that all Values must implement. For now
// it represents any JSON value (object | array | boolean | string | number | null)
type ValueContent interface{}

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
	PropertyValue

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
