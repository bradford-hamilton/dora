// Package lexer TODO
package lexer

import (
	"github.com/bradford-hamilton/dora/pkg/token"
)

// Lexer holds input data and fields that help with scanning.
// It's methods perform lexical analysis/scanning.
type Lexer struct {
	Input        []byte
	char         byte // current char under examination
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	line         int  // line number for better error reporting, etc
}

// New creates and returns a pointer to the Lexer
func New(input string) *Lexer {
	l := &Lexer{Input: []byte(input)}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.Input) {
		// End of input (haven't read anything yet or EOF)
		// 0 is ASCII code for "NUL" character
		l.char = 0
	} else {
		l.char = l.Input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

// NextToken switches through the lexer's current char and creates a new token.
// It then it calls readChar() to advance the lexer and it returns the token
func (l *Lexer) NextToken() token.Token {
	var t token.Token

	if l.isWhitespace() {
		t.Type = token.Whitespace
		t.Line = l.line
		t.Start = l.position
		t.Literal = l.readWhitespace()
		t.End = l.position
		return t
	}

	switch l.char {
	case '/':
		t.Start = l.position
		t.Line = l.line
		l.readChar()
		if l.char != '/' {
			t = newToken(token.Illegal, l.line, 1, 2, l.char)
			return t
		}
		t.Type = token.LineComment
		t.Literal = "/" + l.readLine()
		t.End = l.position
		return t // skip the default readChar
	case '{':
		t = newToken(token.LeftBrace, l.line, l.position, l.position+1, l.char)
	case '}':
		t = newToken(token.RightBrace, l.line, l.position, l.position+1, l.char)
	case '[':
		t = newToken(token.LeftBracket, l.line, l.position, l.position+1, l.char)
	case ']':
		t = newToken(token.RightBracket, l.line, l.position, l.position+1, l.char)
	case ':':
		t = newToken(token.Colon, l.line, l.position, l.position+1, l.char)
	case ',':
		t = newToken(token.Comma, l.line, l.position, l.position+1, l.char)
	case '"', '\'':
		delimiter := l.char
		t.Type = token.String
		t.Literal = l.readString(delimiter)
		t.Line = l.line
		t.Start = l.position
		t.End = l.position + 1
		t.Prefix = string(delimiter)
		t.Suffix = string(delimiter)
	case 0:
		t.Literal = ""
		t.Type = token.EOF
		t.Line = l.line
	default:
		if isLetter(l.char) {
			t.Start = l.position
			ident := l.readIdentifier()
			t.Literal = ident
			t.Line = l.line
			t.End = l.position

			tokenType, err := token.LookupIdentifier(ident)
			if err != nil {
				t.Type = token.Illegal
				return t
			}
			t.Type = tokenType
			t.End = l.position
			return t
		} else if isNumber(l.char) {
			t.Start = l.position
			t.Literal = l.readNumber()
			t.Type = token.Number
			t.Line = l.line
			t.End = l.position
			return t
		}
		t = newToken(token.Illegal, l.line, 1, 2, l.char)
	}

	l.readChar()

	return t
}

func (l *Lexer) isWhitespace() bool {
	return l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r'
}

func (l *Lexer) readWhitespace() string {
	result := ""
	for l.isWhitespace() {
		if l.char == '\n' {
			l.line++
		}
		result += string(l.char)
		l.readChar() // advance
	}
	return result
}

func (l *Lexer) skipWhitespace() {
	for l.isWhitespace() {
		if l.char == '\n' {
			l.line++
		}
		l.readChar()
	}
}

func newToken(tokenType token.Type, line, start, end int, char ...byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(char),
		Line:    line,
		Start:   start,
		End:     end,
	}
}

// readString sets a start position and reads through characters
// When it finds a closing `"`, it stops consuming characters and
// returns the string between the start and end positions.
func (l *Lexer) readString(delimiter byte) string {
	position := l.position + 1
	for {
		prevChar := l.char
		l.readChar()
		if (l.char == delimiter && prevChar != '\\') || l.char == 0 {
			break
		}
	}
	return string(l.Input[position:l.position])
}

// readLine sets a start position and reads through characters
// When it finds a line break, it stops consuming characters and
// returns the string between the start and end positions.
func (l *Lexer) readLine() string {
	position := l.position
	for {
		l.readChar()
		if l.char == 0 {
			break
		}
		if l.char == '\n' {
			l.line++
			l.readChar()
			break
		}
	}
	return string(l.Input[position:l.position])
}

// readNumber sets a start position and reads through characters. When it
// finds a char that isn't a number, it stops consuming characters and
// returns the string between the start and end positions.
func (l *Lexer) readNumber() string {
	position := l.position

	for isNumber(l.char) {
		l.readChar()
	}

	return string(l.Input[position:l.position])
}

func isNumber(char byte) bool {
	return '0' <= char && char <= '9' || char == '.' || char == '-'
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z'
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for isLetter(l.char) {
		l.readChar()
	}

	return string(l.Input[position:l.position])
}
