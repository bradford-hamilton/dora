// Package dora TODO: package docs
package dora

import (
	"github.com/bradford-hamilton/dora/pkg/ast"
	"github.com/bradford-hamilton/dora/pkg/lexer"
	"github.com/bradford-hamilton/dora/pkg/parser"
)

// Client represents a dora client. The client holds things like a copy of the input, the tree (the
// parsed AST representation built with Go types), the user's query & parsed version of the query, and
// a query result. Client exposes public methods which access this underlying data.
type Client struct {
	input       []rune
	tree        *ast.RootNode
	query       []rune
	parsedQuery []queryToken
	result      string
}

// NewFromString takes a string, creates a lexer, creates a parser from the lexer,
// and parses the json into an AST. Methods on the Client give access to private
// data like the AST held inside.
func NewFromString(jsonStr string) (*Client, error) {
	l := lexer.New(jsonStr)
	p := parser.New(l)
	tree, err := p.ParseJSON()
	if err != nil {
		return nil, err
	}
	return &Client{tree: &tree, input: l.Input}, nil
}

// NewFromBytes takes a slice of bytes, converts it to a string, then returns `NewFromString`, passing in the JSON string.
func NewFromBytes(bytes []byte) (*Client, error) {
	return NewFromString(string(bytes))
}

// GetByPath takes a dora query, prepares and validates it, executes the query, and returns the result or an error.
func (c *Client) GetByPath(query string) (string, error) {
	if err := c.prepareQuery(query, c.tree.Type); err != nil {
		return "", err
	}
	if err := c.executeQuery(); err != nil {
		return "", err
	}
	return c.result, nil
}
