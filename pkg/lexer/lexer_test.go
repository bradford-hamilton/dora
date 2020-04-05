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
			"ppu": 0.55,
			"batters": {
				"batter": [{
					"id": "1001",
					"type": "Regular",
					"fun": "true"
				}]
			},
			"topping": [{
				"id": "5001",
				"type": "null",
				"fun": "false"
			}]
		}]
	}
}`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
		expectedLine    int
	}{
		{token.LeftBrace, "{", 0},
		{token.String, "items", 1},
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
