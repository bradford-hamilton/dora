// Package dora TODO: package docs
package dora

import (
	"errors"
	"fmt"

	"github.com/bradford-hamilton/dora/pkg/ast"
	"github.com/bradford-hamilton/dora/pkg/lexer"
	"github.com/bradford-hamilton/dora/pkg/parser"
)

// ErrGetFromNullObj TODO: maybe define specific error types, ex:
// var errAttrNotFound = errors.New("requested attribute not available")

// Client represents a dora client. It holds a private field "program" which holds the AST.
// Methods on this client are the public API for dora.
type Client struct {
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
	return &Client{program: &program}, nil
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

// prepareQuery validates the query root, sets the query on the client struct, and parses the query into parsedQuery of some sort.
// TODO: later decide if this is useful
func (c *Client) prepareQuery(query string, rootNodeType ast.RootNodeType) error {
	if err := validateQueryRoot(query, c.program.Type); err != nil {
		return err
	}
	c.setQuery([]rune(query))
	if err := c.parseQuery(); err != nil {
		return err
	}
	return nil
}

func (c *Client) setQuery(query []rune) {
	c.query = query
}

func (c *Client) parseQuery() error {
	tokens, err := scanQueryTokens(c.query)
	if err != nil {
		return err
	}
	c.parsedQuery = tokens
	return nil
}

func (c *Client) executeQuery() error {
	// TODO: actually execute the query next
	fmt.Println(c.query)
	// parseQuery into parsedQuery type deal?
	// iterate through c.program to fetch user request
	return nil
}

func validateQueryRoot(query string, rootNodeType ast.RootNodeType) error {
	if query[0] != '$' {
		return ErrNoDollarSignRoot
	}

	// The query root after the `$` must be a `.` if the rootNodeType is an object
	validObjQueryRoot := query[1] == '.'
	if rootNodeType == ast.ObjectRoot && !validObjQueryRoot {
		return ErrWrongObjectRootSelector
	}

	// The query root after the `$` must be a `[` if the rootNodeType is an array
	validArrayQueryRoot := query[1] == '['
	if rootNodeType == ast.ArrayRoot && !validArrayQueryRoot {
		return ErrWrongArrayRootSelector
	}

	return nil
}

// ErrNoDollarSignRoot is used for telling the user the very first character must be a `$`
var ErrNoDollarSignRoot = errors.New("Incorrect syntax, query must start with `$` representing the root object or array")

// ErrWrongObjectRootSelector is used for telling the user their JSON root is an object and the selector found was not a `.`
var ErrWrongObjectRootSelector = errors.New("Incorrect syntax. Your root JSON type is an object. Therefore, path queries must begin by selecting a `key` from your root object. Ex: `$.keyOnRootObject` or `$[\"keyOnRootObject\"]`")

// ErrWrongArrayRootSelector is used for telling the user their JSON root is an array and the selector found was not a `[`
var ErrWrongArrayRootSelector = errors.New("Incorrect syntax. Your root JSON type is an array. Therefore, path queries must begin by selecting an item by index on the root array. Ex: `$[0]` or `$[1]`")

func errSelectorSytax(operator string) error {
	return fmt.Errorf(
		"Error parsing query, expected either a `.` for selections on an object or a `[` for selections on an array. Got: %s",
		operator,
	)
}
