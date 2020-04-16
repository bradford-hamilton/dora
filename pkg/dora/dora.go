// Package dora TODO: package docs
package dora

import (
	"github.com/bradford-hamilton/dora/pkg/ast"
	"github.com/bradford-hamilton/dora/pkg/lexer"
	"github.com/bradford-hamilton/dora/pkg/parser"
)

// Client represents a dora client. It exposes public methods which access to the underlying data
type Client struct {
	input       []rune
	program     *ast.RootNode
	query       []rune
	parsedQuery []queryToken
	result      string
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
	return &Client{program: &program, input: l.Input}, nil
}

// NewFromBytes takes a slice of bytes, converts it to a string, then returns `NewFromString`, passing in the JSON string.
func NewFromBytes(bytes []byte) (*Client, error) {
	return NewFromString(string(bytes))
}

// GetByPath thinking about what methods to add to the client to interact with the program
func (c *Client) GetByPath(query string) (string, error) {
	if err := c.prepareQuery(query, c.program.Type); err != nil {
		return "", err
	}
	if err := c.executeQuery(); err != nil {
		return "", err
	}
	return c.result, nil
}
