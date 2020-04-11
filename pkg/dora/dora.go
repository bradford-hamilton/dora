// Package dora TODO: package docs
package dora

import (
	"errors"
	"fmt"

	"github.com/bradford-hamilton/dora/pkg/ast"
	"github.com/bradford-hamilton/dora/pkg/lexer"
	"github.com/bradford-hamilton/dora/pkg/parser"
)

// Client represents a dora client. It holds a private field "program" which holds the AST.
// Methods on this client are the public API for dora.
type Client struct {
	program *ast.RootNode
	query   string
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
	str := string(bytes)
	return NewFromString(str)
}

// GetByFullPath thinking about what methods to add to the client to interact with the program
func (c *Client) GetByFullPath(query string) (string, error) {
	if err := validateQueryRoot(query, c.program.Type); err != nil {
		return "", err
	}
	c.query = query

	result, err := c.executeQuery()
	if err != nil {
		return "", err
	}

	return result, nil
}

func (c *Client) executeQuery() (string, error) {
	fmt.Println(c.query)
	// parseQuery into steps type deal?
	// iterate through c.program to fetch user request

	return "", nil
}

func validateQueryRoot(query string, rootNodeType ast.RootNodeType) error {
	if query[0] != '$' {
		return errors.New(
			"Incorrect syntax, query must start with `$` representing the root object or array",
		)
	}

	validObjQueryRoot := query[1] == '.' || query[1] == '['
	if rootNodeType == ast.ObjectRoot && !validObjQueryRoot {
		return errors.New(
			"Incorrect syntax. Your root JSON type is an object. Therefore, path queries must begin by selecting a `key` from your root object. Ex: `$.keyOnRootObject` or `$[\"keyOnRootObject\"]`",
		)
	}

	validArrayQueryRoot := query[1] == '['
	if rootNodeType == ast.ArrayRoot && !validArrayQueryRoot {
		return errors.New(
			"Incorrect syntax. Your root JSON type is an array. Therefore, path queries must begin by selecting an item by index on the root array. Ex: `$[0]` or `$[1]`",
		)
	}

	return nil
}

/*
	JSON Object:
	{
		"name": "bradford",
		"array": ["some", "values"]
	}

	- Query must start with $.
	$.name == "bradford"
	$.array == ["array", "values"]
	$.array[0] == "some"
	$.array[1] == "values"
	$.array[2] == error

	---

	JSON Array:
	[
		"some",
		"values",
		{ "objKey": "objValue" }
	]

	- Query must start with $.
	$[0] == "some"
	$[1] == "values"
	$[2] == { "objKey": "objValue" }
	$[2].objKey == "objValue"
	$[2]["objKey"] == "objValue"

*/
