package parser

import (
	"testing"

	"github.com/bradford-hamilton/dora/pkg/ast"
	"github.com/bradford-hamilton/dora/pkg/lexer"
)

func TestParsingJSONObjectChildren(t *testing.T) {
	tests := []struct {
		input       string
		childrenLen int
	}{
		{"{\"key0\": \"value0\"}", 1},
		{"{\"key1\": \"value1\", \"key2\": \"value2\"}", 2},
		{"{\"key3\": [\"value3\", \"value4\"]}", 1},
		{"{\"key4\": [\"value5\", {\"key5\": \"value6\"}]}", 1},
		{"{\"key5\":\" value7\", \"key6\": \"value7\"}", 2},
		{"{\"key5\":\" value7\", \"key6\": \"value7\", \"key7\": \"value8\"}", 3},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("Failed to parse program. Error: %v", err)
		}

		rv := *program.RootValue
		val := rv.(ast.Object)

		checkParserErrors(t, p)

		if len(val.Children) != tt.childrenLen {
			t.Fatalf("The length of the children does not contain 1 statement. Got: %d", len(val.Children))
		}
	}
}

func TestParsingJSONArrayChildren(t *testing.T) {
	tests := []struct {
		input       string
		childrenLen int
	}{
		{"[\"item1\"]", 1},
		{"[\"item2\", \"item3\"]", 2},
		{"[\"item4\", \"item5\", \"item6\"]", 3},
		{"[\"item7\", \"item8\", {\"key1\": \"value1\"}]", 3},
		{"[\"item9\", \"item10\", {\"key1\": \"value1\"}, \"item11\"]", 4},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("Failed to parse program. Error: %v", err)
		}

		rv := *program.RootValue
		val := rv.(ast.Array)

		checkParserErrors(t, p)

		if len(val.Children) != tt.childrenLen {
			t.Fatalf("The length of the children does not contain 1 statement. Got: %d", len(val.Children))
		}
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
