// Package parsejson TODO: package docs
package parsejson

import (
	"github.com/bradford-hamilton/parsejson/pkg/ast"
	"github.com/bradford-hamilton/parsejson/pkg/lexer"
	"github.com/bradford-hamilton/parsejson/pkg/parser"
)

// Client TODO
type Client struct {
	program *ast.RootNode
}

// NewFromBytes takes a slice of bytes, creates a lexer, creates a parser from the lexer,
// and parses the program into an AST. Methods on the Client give access to private
// AST held inside
func NewFromBytes(bytes []byte) (*Client, error) {
	l := lexer.New(string(bytes))
	p := parser.New(l)
	program, err := p.ParseProgram()
	if err != nil {
		return nil, err
	}
	return &Client{&program}, nil
}

// NewFromString takes a string, creates a lexer, creates a parser from the lexer,
// and parses the program into an AST. Methods on the Client give access to private
// AST held inside
func NewFromString(jsonStr string) (*Client, error) {
	l := lexer.New(jsonStr)
	p := parser.New(l)
	program, err := p.ParseProgram()
	if err != nil {
		return nil, err
	}
	return &Client{&program}, nil
}

// Start thinking about what methods to add to the client to interact with the program
// func (p *Client) Get() {}
