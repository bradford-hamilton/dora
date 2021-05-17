package lexer

import (
	"fmt"
	"testing"

	"github.com/bradford-hamilton/dora/pkg/token"
)

func TestNextToken_WithSingleLineComments(t *testing.T) {
	input := `// Initial comment
{
	"name" : "Stuart", // test comment
}
// ending comment
`

	tests := []token.Token{
		{Type: token.LineComment, Literal: "// Initial comment\n", Line: 0},
		{Type: token.LeftBrace, Literal: "{", Line: 1},
		{Type: token.Whitespace, Literal: "\n\t", Line: 1},
		{Type: token.String, Literal: "name", Line: 2, Prefix: `"`, Suffix: `"`},
		{Type: token.Whitespace, Literal: " ", Line: 2},
		{Type: token.Colon, Literal: ":", Line: 2},
		{Type: token.Whitespace, Literal: " ", Line: 2},
		{Type: token.String, Literal: "Stuart", Line: 2, Prefix: `"`, Suffix: `"`},
		{Type: token.Comma, Literal: ",", Line: 2},
		{Type: token.Whitespace, Literal: " ", Line: 2},
		{Type: token.LineComment, Literal: "// test comment\n", Line: 2},
		{Type: token.RightBrace, Literal: "}", Line: 3},
		{Type: token.Whitespace, Literal: "\n", Line: 3},
		{Type: token.LineComment, Literal: "// ending comment\n", Line: 4},
		{Type: token.EOF, Literal: "", Line: 5},
	}

	l := New(input)

	assertLexerMatches(t, l, tests)
}

func TestNextToken_WithSingleQuoteString(t *testing.T) {
	input := `'"name"''`

	tests := []token.Token{
		{Type: token.String, Literal: `"name"`, Line: 0, Prefix: `'`, Suffix: `'`},
	}

	l := New(input)

	assertLexerMatches(t, l, tests)
}

func TestNextToken(t *testing.T) {
	input := `{
	"items": {
		"item": [{
			"id": "0001",
			"type": "donut",
			"name": "Cake",
			"cpu": 55,
			"batters": {
				"batter": [{
					"id": false,
					"name": null,
					"fun": true
				}]
			},
			"names": ["catstack", "lampcat", "langlang"]
		}]
	},
	"version": 0.1,
	"number": 11.4,
	"negativeNum": -5,
	"escapeString": "I'm some \"string\" thats escaped"
}`

	tests := []token.Token{
		{Type: token.LeftBrace, Literal: "{", Line: 0},
		{Type: token.Whitespace, Literal: "\n\t", Line: 0},
		{Type: token.String, Literal: "items", Line: 1, Prefix: `"`, Suffix: `"`},
		{Type: token.Colon, Literal: ":", Line: 1},
		{Type: token.Whitespace, Literal: " ", Line: 1},
		{Type: token.LeftBrace, Literal: "{", Line: 1},
		{Type: token.Whitespace, Literal: "\n\t\t", Line: 1},
		{Type: token.String, Literal: "item", Line: 2, Prefix: `"`, Suffix: `"`},
		{Type: token.Colon, Literal: ":", Line: 2},
		{Type: token.Whitespace, Literal: " ", Line: 2},
		{Type: token.LeftBracket, Literal: "[", Line: 2},
		{Type: token.LeftBrace, Literal: "{", Line: 2},
		{Type: token.Whitespace, Literal: "\n\t\t\t", Line: 2},
		{Type: token.String, Literal: "id", Line: 3, Prefix: `"`, Suffix: `"`},
		{Type: token.Colon, Literal: ":", Line: 3},
		{Type: token.Whitespace, Literal: " ", Line: 3},
		{Type: token.String, Literal: "0001", Line: 3, Prefix: `"`, Suffix: `"`},
		{Type: token.Comma, Literal: ",", Line: 3},
		{Type: token.Whitespace, Literal: "\n\t\t\t", Line: 3},
		{Type: token.String, Literal: "type", Line: 4, Prefix: `"`, Suffix: `"`},
		{Type: token.Colon, Literal: ":", Line: 4},
		{Type: token.Whitespace, Literal: " ", Line: 4},
		{Type: token.String, Literal: "donut", Line: 4, Prefix: `"`, Suffix: `"`},
		{Type: token.Comma, Literal: ",", Line: 4},
		{Type: token.Whitespace, Literal: "\n\t\t\t", Line: 4},
		{Type: token.String, Literal: "name", Line: 5, Prefix: `"`, Suffix: `"`},
		{Type: token.Colon, Literal: ":", Line: 5},
		{Type: token.Whitespace, Literal: " ", Line: 5},
		{Type: token.String, Literal: "Cake", Line: 5, Prefix: `"`, Suffix: `"`},
		{Type: token.Comma, Literal: ",", Line: 5},
		{Type: token.Whitespace, Literal: "\n\t\t\t", Line: 5},
		{Type: token.String, Literal: "cpu", Line: 6, Prefix: `"`, Suffix: `"`},
		{Type: token.Colon, Literal: ":", Line: 6},
		{Type: token.Whitespace, Literal: " ", Line: 6},
		{Type: token.Number, Literal: "55", Line: 6},
		{Type: token.Comma, Literal: ",", Line: 6},
		{Type: token.Whitespace, Literal: "\n\t\t\t", Line: 6},
		{Type: token.String, Literal: "batters", Line: 7, Prefix: `"`, Suffix: `"`},
		{Type: token.Colon, Literal: ":", Line: 7},
		{Type: token.Whitespace, Literal: " ", Line: 7},
		{Type: token.LeftBrace, Literal: "{", Line: 7},
		{Type: token.Whitespace, Literal: "\n\t\t\t\t", Line: 7},
		{Type: token.String, Literal: "batter", Line: 8, Prefix: `"`, Suffix: `"`},
		{Type: token.Colon, Literal: ":", Line: 8},
		{Type: token.Whitespace, Literal: " ", Line: 8},
		{Type: token.LeftBracket, Literal: "[", Line: 8},
		{Type: token.LeftBrace, Literal: "{", Line: 8},
		{Type: token.Whitespace, Literal: "\n\t\t\t\t\t", Line: 8},
		{Type: token.String, Literal: "id", Line: 9, Prefix: `"`, Suffix: `"`},
		{Type: token.Colon, Literal: ":", Line: 9},
		{Type: token.Whitespace, Literal: " ", Line: 9},
		{Type: token.False, Literal: "false", Line: 9},
		{Type: token.Comma, Literal: ",", Line: 9},
		{Type: token.Whitespace, Literal: "\n\t\t\t\t\t", Line: 9},
		{Type: token.String, Literal: "name", Line: 10, Prefix: `"`, Suffix: `"`},
		{Type: token.Colon, Literal: ":", Line: 10},
		{Type: token.Whitespace, Literal: " ", Line: 10},
		{Type: token.Null, Literal: "null", Line: 10},
		{Type: token.Comma, Literal: ",", Line: 10},
		{Type: token.Whitespace, Literal: "\n\t\t\t\t\t", Line: 10},
		{Type: token.String, Literal: "fun", Line: 11, Prefix: `"`, Suffix: `"`},
		{Type: token.Colon, Literal: ":", Line: 11},
		{Type: token.Whitespace, Literal: " ", Line: 11},
		{Type: token.True, Literal: "true", Line: 11},
		{Type: token.Whitespace, Literal: "\n\t\t\t\t", Line: 11},
		{Type: token.RightBrace, Literal: "}", Line: 12},
		{Type: token.RightBracket, Literal: "]", Line: 12},
		{Type: token.Whitespace, Literal: "\n\t\t\t", Line: 12},
		{Type: token.RightBrace, Literal: "}", Line: 13},
		{Type: token.Comma, Literal: ",", Line: 13},
		{Type: token.Whitespace, Literal: "\n\t\t\t", Line: 13},
		{Type: token.String, Literal: "names", Line: 14, Prefix: `"`, Suffix: `"`},
		{Type: token.Colon, Literal: ":", Line: 14},
		{Type: token.Whitespace, Literal: " ", Line: 14},
		{Type: token.LeftBracket, Literal: "[", Line: 14},
		{Type: token.String, Literal: "catstack", Line: 14, Prefix: `"`, Suffix: `"`},
		{Type: token.Comma, Literal: ",", Line: 14},
		{Type: token.Whitespace, Literal: " ", Line: 14},
		{Type: token.String, Literal: "lampcat", Line: 14, Prefix: `"`, Suffix: `"`},
		{Type: token.Comma, Literal: ",", Line: 14},
		{Type: token.Whitespace, Literal: " ", Line: 14},
		{Type: token.String, Literal: "langlang", Line: 14, Prefix: `"`, Suffix: `"`},
		{Type: token.RightBracket, Literal: "]", Line: 14},
		{Type: token.Whitespace, Literal: "\n\t\t", Line: 14},
		{Type: token.RightBrace, Literal: "}", Line: 15},
		{Type: token.RightBracket, Literal: "]", Line: 15},
		{Type: token.Whitespace, Literal: "\n\t", Line: 15},
		{Type: token.RightBrace, Literal: "}", Line: 16},
		{Type: token.Comma, Literal: ",", Line: 16},
		{Type: token.Whitespace, Literal: "\n\t", Line: 16},
		{Type: token.String, Literal: "version", Line: 17, Prefix: `"`, Suffix: `"`},
		{Type: token.Colon, Literal: ":", Line: 17},
		{Type: token.Whitespace, Literal: " ", Line: 17},
		{Type: token.Number, Literal: "0.1", Line: 17},
		{Type: token.Comma, Literal: ",", Line: 17},
		{Type: token.Whitespace, Literal: "\n\t", Line: 17},
		{Type: token.String, Literal: "number", Line: 18, Prefix: `"`, Suffix: `"`},
		{Type: token.Colon, Literal: ":", Line: 18},
		{Type: token.Whitespace, Literal: " ", Line: 18},
		{Type: token.Number, Literal: "11.4", Line: 18},
		{Type: token.Comma, Literal: ",", Line: 18},
		{Type: token.Whitespace, Literal: "\n\t", Line: 18},
		{Type: token.String, Literal: "negativeNum", Line: 19, Prefix: `"`, Suffix: `"`},
		{Type: token.Colon, Literal: ":", Line: 19},
		{Type: token.Whitespace, Literal: " ", Line: 19},
		{Type: token.Number, Literal: "-5", Line: 19},
		{Type: token.Comma, Literal: ",", Line: 19},
		{Type: token.Whitespace, Literal: "\n\t", Line: 19},
		{Type: token.String, Literal: "escapeString", Line: 20, Prefix: `"`, Suffix: `"`},
		{Type: token.Colon, Literal: ":", Line: 20},
		{Type: token.Whitespace, Literal: " ", Line: 20},
		{Type: token.String, Literal: "I'm some \\\"string\\\" thats escaped", Line: 20, Prefix: `"`, Suffix: `"`},
		{Type: token.Whitespace, Literal: "\n", Line: 20},
		{Type: token.RightBrace, Literal: "}", Line: 21},
		{Type: token.EOF, Literal: "", Line: 21},
	}

	l := New(input)

	assertLexerMatches(t, l, tests)
}

func assertLexerMatches(t *testing.T, l *Lexer, tests []token.Token) {
	for i, expectedToken := range tests {
		token := l.NextToken()

		if token.Type != expectedToken.Type {
			t.Fatalf("tests[%d] - tokentype wrong. Expected: %s, Got: %s", i, formatTokenOutputString(expectedToken), formatTokenOutputString(token))
		}
		if token.Literal != expectedToken.Literal {
			t.Fatalf("tests[%d] - literal wrong. Expected: %s, Got: %s", i, formatTokenOutputString(expectedToken), formatTokenOutputString(token))
		}
		if token.Line != expectedToken.Line {
			t.Fatalf("tests[%d] - line wrong. Expected: %s, Got: %s", i, formatTokenOutputString(expectedToken), formatTokenOutputString(token))
		}
		if token.Prefix != expectedToken.Prefix {
			t.Fatalf("tests[%d] - prefix wrong. Expected: %s, Got: %s", i, formatTokenOutputString(expectedToken), formatTokenOutputString(token))
		}
		if token.Suffix != expectedToken.Suffix {
			t.Fatalf("tests[%d] - suffix wrong. Expected: %s, Got: %s", i, formatTokenOutputString(expectedToken), formatTokenOutputString(token))
		}
	}
}

func formatTokenOutputString(t token.Token) string {
	result := fmt.Sprintf("Type:%q; Literal:%q; Line:%d", t.Type, t.Literal, t.Line)
	if t.Prefix != "" {
		result += fmt.Sprintf("; Prefix:%q", t.Prefix)
	}
	if t.Suffix != "" {
		result += fmt.Sprintf("; Suffix:%q", t.Suffix)
	}
	return result
}
