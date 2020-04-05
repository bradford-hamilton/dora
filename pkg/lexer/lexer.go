package lexer

import (
	"github.com/bradford-hamilton/parsejson/pkg/token"
)

// Lexer performs lexical analysis/scanning of the JSON
type Lexer struct {
	input        []rune
	char         rune // current char under examination
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	line         int  // line number for better error reporting, etc
}

// New creates and returns a pointer to the Lexer
func New(input string) *Lexer {
	l := &Lexer{input: []rune(input)}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		// End of input (haven't read anything yet or EOF)
		// 0 is ASCII code for "NUL" character
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

// NextToken switches through the lexer's current char and creates a new token.
// It then it calls readChar() to advance the lexer and it returns the token
func (l *Lexer) NextToken() token.Token {
	var t token.Token

	l.skipWhitespace()

	switch l.char {
	case '{':
		t = newToken(token.LeftBrace, l.line, l.char)
	case '}':
		t = newToken(token.RightBrace, l.line, l.char)
	case '[':
		t = newToken(token.LeftBracket, l.line, l.char)
	case ']':
		t = newToken(token.RightBracket, l.line, l.char)
	case ':':
		t = newToken(token.Colon, l.line, l.char)
	case ',':
		t = newToken(token.Comma, l.line, l.char)
	case '"':
		t.Type = token.String
		t.Literal = l.readString()
		t.Line = l.line

		if t.Literal == "true" {
			t.Type = token.True
		}
		if t.Literal == "false" {
			t.Type = token.False
		}
	case 0:
		t.Literal = ""
		t.Type = token.EOF
		t.Line = l.line
	default:
		if isInteger(l.char) {
			t.Literal = l.readInteger()
			t.Type = token.Integer
			t.Line = l.line
			return t
		}
		t = newToken(token.Illegal, l.line, l.char)
	}

	l.readChar()

	return t
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' || l.char == '\f' || l.char == '\b' {
		if l.char == '\n' {
			l.line++
		}
		l.readChar()
	}
}

func newToken(tokenType token.Type, line int, char ...rune) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(char),
		Line:    line,
	}
}

// readString sets a start position and reads through characters
// When it finds a closing `"`, it stops consuming characters and
// returns the string between the start and end positions.
func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.char == '"' || l.char == 0 {
			break
		}
	}
	return string(l.input[position:l.position])
}

// readInteger sets a start position and reads through characters. When it
// finds a char that isn't an integer, it stops consuming characters and
// returns the string between the start and end positions.
func (l *Lexer) readInteger() string {
	position := l.position

	for isInteger(l.char) {
		l.readChar()
	}

	return string(l.input[position:l.position])
}

func isInteger(char rune) bool {
	return '0' <= char && char <= '9'
}
