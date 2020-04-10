package parser

import (
	"fmt"
	"testing"

	"github.com/bradford-hamilton/parsejson/pkg/lexer"
)

func TestParsingJSON(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"{\"key\": \"value\"}"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("Failed to parse program. Error: %v", err)
		}
		fmt.Println(program)
		// t := program.RootValue.(*ast.Object)
		checkParserErrors(t, p)

		// if len(program.RootValue) != 1 {
		// 	t.Fatalf("program.Statements does not contain 1 statement. Got: %d", len(program.Statements))
		// }
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("Parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}
	t.FailNow()
}
