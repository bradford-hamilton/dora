// Package token TODO
package token

import (
	"fmt"
)

// All the different tokens for supporting JSON
const (
	// Token/character we don't know about
	Illegal = "ILLEGAL"

	// End of file
	EOF = "EOF"

	// Literals
	String  = "STRING"
	Integer = "INTEGER"
	Float   = "FLOAT"

	// The six structural tokens
	LeftBrace    = "{"
	RightBrace   = "}"
	LeftBracket  = "["
	RightBracket = "]"
	Comma        = ","
	Colon        = ":"

	// Values
	True  = "TRUE"
	False = "FALSE"
	Null  = "NULL"
)

// Type is a type alias for a string
type Type string

// Token is a struct representing a JSON token - holds a Type and a literal
type Token struct {
	Type    Type
	Literal string
	Line    int
}

var validJSONIdentifiers = map[string]Type{
	"true":  True,
	"false": False,
	"null":  Null,
}

// LookupIdentifier checks our validJSONIdentifiers map for the scanned identifier. If it finds one,
// the identifier's token type is returned. If not found, an error is returned
func LookupIdentifier(identifier string) (Type, error) {
	if token, ok := validJSONIdentifiers[identifier]; ok {
		return token, nil
	}
	return "", fmt.Errorf("Expected a valid JSON identifier. Found: %s", identifier)
}

// https://www.ecma-international.org/publications/files/ECMA-ST/ECMA-404.pdf

// [ U+005B left square bracket
// { U+007B left curly bracket
// ] U+005D right square bracket
// } U+007D right curly bracket
// : U+003A colon
// , U+002C comma
// These are the three literal name tokens:
// true U+0074 U+0072 U+0075 U+0065
// false U+0066 U+0061 U+006C U+0073 U+0065
// null U+006E U+0075 U+006C U+006C

// \" represents the quotation mark character (U+0022).
// \\ represents the reverse solidus character (U+005C).
// \/ represents the solidus character (U+002F).
// \b represents the backspace character (U+0008).
// \f represents the form feed character (U+000C).
// \n represents the line feed character (U+000A).
// \r represents the carriage return character (U+000D).
// \t represents the character tabulation character (U+0009).
