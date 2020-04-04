package token

// All the different tokens for supporting JSON
const (
	// Token/character we don't know about
	Illegal = "ILLEGAL"

	// End of file
	EOF = "EOF"

	// The six structural tokens
	LeftBrace    = "{"
	RightBrace   = "}"
	LeftBracket  = "["
	RightBracket = "]"
	Comma        = ","
	Colon        = ":"

	// The three literal name tokens
	True  = "true"
	False = "false"
	Null  = "null"
)

// Type is a type alias for a string
type Type string

// Token is a struct representing a JSON token - holds a Type and a literal
type Token struct {
	Type    Type
	Literal string
	Line    int
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
