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

	tests := [...]struct {
		expectedType    token.Type
		expectedLiteral string
		expectedLine    int
	}{
		{token.LineComment, "// Initial comment\n", 0},
		{token.LeftBrace, "{", 1},
		{token.Whitespace, "\n\t", 1},
		{token.String, "name", 2},
		{token.Whitespace, " ", 2},
		{token.Colon, ":", 2},
		{token.Whitespace, " ", 2},
		{token.String, "Stuart", 2},
		{token.Comma, ",", 2},
		{token.Whitespace, " ", 2},
		{token.LineComment, "// test comment\n", 2},
		{token.RightBrace, "}", 3},
		{token.Whitespace, "\n", 3},
		{token.LineComment, "// ending comment\n", 4},
		{token.EOF, "", 5},
	}

	l := New(input)

	for i, tt := range tests {
		token := l.NextToken()

		expectedText := fmt.Sprintf("%q %q %d", tt.expectedType, tt.expectedLiteral, tt.expectedLine)
		actualText := fmt.Sprintf("%q %q %d", token.Type, token.Literal, token.Line)

		if token.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. Expected: %s, Got: %s", i, expectedText, actualText)
		}
		if token.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. Expected: %s, Got: %s", i, expectedText, actualText)
		}
		if token.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong. Expected: %s, Got: %s", i, expectedText, actualText)
		}
	}
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

	tests := [...]struct {
		expectedType    token.Type
		expectedLiteral string
		expectedLine    int
	}{
		{token.LeftBrace, "{", 0},
		{token.Whitespace, "\n\t", 0},
		{token.String, "items", 1},
		{token.Colon, ":", 1},
		{token.Whitespace, " ", 1},
		{token.LeftBrace, "{", 1},
		{token.Whitespace, "\n\t\t", 1},
		{token.String, "item", 2},
		{token.Colon, ":", 2},
		{token.Whitespace, " ", 2},
		{token.LeftBracket, "[", 2},
		{token.LeftBrace, "{", 2},
		{token.Whitespace, "\n\t\t\t", 2},
		{token.String, "id", 3},
		{token.Colon, ":", 3},
		{token.Whitespace, " ", 3},
		{token.String, "0001", 3},
		{token.Comma, ",", 3},
		{token.Whitespace, "\n\t\t\t", 3},
		{token.String, "type", 4},
		{token.Colon, ":", 4},
		{token.Whitespace, " ", 4},
		{token.String, "donut", 4},
		{token.Comma, ",", 4},
		{token.Whitespace, "\n\t\t\t", 4},
		{token.String, "name", 5},
		{token.Colon, ":", 5},
		{token.Whitespace, " ", 5},
		{token.String, "Cake", 5},
		{token.Comma, ",", 5},
		{token.Whitespace, "\n\t\t\t", 5},
		{token.String, "cpu", 6},
		{token.Colon, ":", 6},
		{token.Whitespace, " ", 6},
		{token.Number, "55", 6},
		{token.Comma, ",", 6},
		{token.Whitespace, "\n\t\t\t", 6},
		{token.String, "batters", 7},
		{token.Colon, ":", 7},
		{token.Whitespace, " ", 7},
		{token.LeftBrace, "{", 7},
		{token.Whitespace, "\n\t\t\t\t", 7},
		{token.String, "batter", 8},
		{token.Colon, ":", 8},
		{token.Whitespace, " ", 8},
		{token.LeftBracket, "[", 8},
		{token.LeftBrace, "{", 8},
		{token.Whitespace, "\n\t\t\t\t\t", 8},
		{token.String, "id", 9},
		{token.Colon, ":", 9},
		{token.Whitespace, " ", 9},
		{token.False, "false", 9},
		{token.Comma, ",", 9},
		{token.Whitespace, "\n\t\t\t\t\t", 9},
		{token.String, "name", 10},
		{token.Colon, ":", 10},
		{token.Whitespace, " ", 10},
		{token.Null, "null", 10},
		{token.Comma, ",", 10},
		{token.Whitespace, "\n\t\t\t\t\t", 10},
		{token.String, "fun", 11},
		{token.Colon, ":", 11},
		{token.Whitespace, " ", 11},
		{token.True, "true", 11},
		{token.Whitespace, "\n\t\t\t\t", 11},
		{token.RightBrace, "}", 12},
		{token.RightBracket, "]", 12},
		{token.Whitespace, "\n\t\t\t", 12},
		{token.RightBrace, "}", 13},
		{token.Comma, ",", 13},
		{token.Whitespace, "\n\t\t\t", 13},
		{token.String, "names", 14},
		{token.Colon, ":", 14},
		{token.Whitespace, " ", 14},
		{token.LeftBracket, "[", 14},
		{token.String, "catstack", 14},
		{token.Comma, ",", 14},
		{token.Whitespace, " ", 14},
		{token.String, "lampcat", 14},
		{token.Comma, ",", 14},
		{token.Whitespace, " ", 14},
		{token.String, "langlang", 14},
		{token.RightBracket, "]", 14},
		{token.Whitespace, "\n\t\t", 14},
		{token.RightBrace, "}", 15},
		{token.RightBracket, "]", 15},
		{token.Whitespace, "\n\t", 15},
		{token.RightBrace, "}", 16},
		{token.Comma, ",", 16},
		{token.Whitespace, "\n\t", 16},
		{token.String, "version", 17},
		{token.Colon, ":", 17},
		{token.Whitespace, " ", 17},
		{token.Number, "0.1", 17},
		{token.Comma, ",", 17},
		{token.Whitespace, "\n\t", 17},
		{token.String, "number", 18},
		{token.Colon, ":", 18},
		{token.Whitespace, " ", 18},
		{token.Number, "11.4", 18},
		{token.Comma, ",", 18},
		{token.Whitespace, "\n\t", 18},
		{token.String, "negativeNum", 19},
		{token.Colon, ":", 19},
		{token.Whitespace, " ", 19},
		{token.Number, "-5", 19},
		{token.Comma, ",", 19},
		{token.Whitespace, "\n\t", 19},
		{token.String, "escapeString", 20},
		{token.Colon, ":", 20},
		{token.Whitespace, " ", 20},
		{token.String, "I'm some \\\"string\\\" thats escaped", 20},
		{token.Whitespace, "\n", 20},
		{token.RightBrace, "}", 21},
		{token.EOF, "", 21},
	}

	l := New(input)

	for i, tt := range tests {
		token := l.NextToken()

		expectedText := fmt.Sprintf("%q %q %d", tt.expectedType, tt.expectedLiteral, tt.expectedLine)
		actualText := fmt.Sprintf("%q %q %d", token.Type, token.Literal, token.Line)

		if token.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. Expected: %s, Got: %s", i, expectedText, actualText)
		}
		if token.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. Expected: %s, Got: %s", i, expectedText, actualText)
		}
		if token.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong. Expected: %s, Got: %s", i, expectedText, actualText)
		}
	}
}
