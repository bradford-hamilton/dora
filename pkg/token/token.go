// Package token TODO
package token

import (
	"fmt"
)

// All the different tokens for supporting JSON
const (
	// Token/character we don't know about
	Illegal Type = "ILLEGAL"

	// End of file
	EOF Type = "EOF"

	// Literals
	String Type = "STRING"
	Number Type = "NUMBER"

	// The six structural tokens
	LeftBrace    Type = "{"
	RightBrace   Type = "}"
	LeftBracket  Type = "["
	RightBracket Type = "]"
	Comma        Type = ","
	Colon        Type = ":"

	// Values
	True  Type = "TRUE"
	False Type = "FALSE"
	Null  Type = "NULL"
)

// Type is a type alias for a string
type Type string

// Token is a struct representing a JSON token - It holds information like its Type and Literal, as well
// as Start, End, and Line fields. Line is used for better error handling, while Start and End are used
// to return objects/arrays from querys.
type Token struct {
	Type    Type
	Literal string
	Line    int
	Start   int
	End     int
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

// TODO: Come back to this for full robust JSON support
var escapes = map[rune]int{
	'"':  0, // Quotation mask
	'\\': 1, // Reverse solidus
	'/':  2, // Solidus
	'b':  3, // Backspace
	'f':  4, // Form feed
	'n':  5, // New line
	'r':  6, // Carriage return
	't':  7, // Horizontal tab
	'u':  8, // 4 hexadecimal digits
}

// TODO: Come back to this for full robust JSON support
var escapeChars = map[string]string{
	"b": "\b", // Backspace
	"f": "\f", // Form feed
	"n": "\n", // New line
	"r": "\r", // Carriage return
	"t": "\t", // Horizontal tab
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
