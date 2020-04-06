package lexer

import (
	"fmt"
	"testing"

	"github.com/bradford-hamilton/parsejson/pkg/token"
)

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
	"number": 11.4
}`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
		expectedLine    int
	}{
		{token.LeftBrace, "{", 0},
		{token.String, "items", 1},
		{token.Colon, ":", 1},
		{token.LeftBrace, "{", 1},
		{token.String, "item", 2},
		{token.Colon, ":", 2},
		{token.LeftBracket, "[", 2},
		{token.LeftBrace, "{", 2},
		{token.String, "id", 3},
		{token.Colon, ":", 3},
		{token.String, "0001", 3},
		{token.Comma, ",", 3},
		{token.String, "type", 4},
		{token.Colon, ":", 4},
		{token.String, "donut", 4},
		{token.Comma, ",", 4},
		{token.String, "name", 5},
		{token.Colon, ":", 5},
		{token.String, "Cake", 5},
		{token.Comma, ",", 5},
		{token.String, "cpu", 6},
		{token.Colon, ":", 6},
		{token.Integer, "55", 6},
		{token.Comma, ",", 6},
		{token.String, "batters", 7},
		{token.Colon, ":", 7},
		{token.LeftBrace, "{", 7},
		{token.String, "batter", 8},
		{token.Colon, ":", 8},
		{token.LeftBracket, "[", 8},
		{token.LeftBrace, "{", 8},
		{token.String, "id", 9},
		{token.Colon, ":", 9},
		{token.False, "false", 9},
		{token.Comma, ",", 9},
		{token.String, "name", 10},
		{token.Colon, ":", 10},
		{token.Null, "null", 10},
		{token.Comma, ",", 10},
		{token.String, "fun", 11},
		{token.Colon, ":", 11},
		{token.True, "true", 11},
		{token.RightBrace, "}", 12},
		{token.RightBracket, "]", 12},
		{token.RightBrace, "}", 13},
		{token.Comma, ",", 13},
		{token.String, "names", 14},
		{token.Colon, ":", 14},
		{token.LeftBracket, "[", 14},
		{token.String, "catstack", 14},
		{token.Comma, ",", 14},
		{token.String, "lampcat", 14},
		{token.Comma, ",", 14},
		{token.String, "langlang", 14},
		{token.RightBracket, "]", 14},
		{token.RightBrace, "}", 15},
		{token.RightBracket, "]", 15},
		{token.RightBrace, "}", 16},
		{token.Comma, ",", 16},
		{token.String, "version", 17},
		{token.Colon, ":", 17},
		// {token.EOF, "", 17},
	}

	l := New(input)

	for i, tt := range tests {
		token := l.NextToken()
		fmt.Printf("Line: %d: Expected: %d", token.Line, tt.expectedLine)
		if token.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. Expected: %q, Got: %q", i, tt.expectedType, token.Type)
		}

		if token.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. Expected: %q, Got: %q", i, tt.expectedLiteral, token.Literal)
		}

		if token.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong. Expected: %q, Got: %q", i, tt.expectedLine, token.Line)
		}
	}
}
