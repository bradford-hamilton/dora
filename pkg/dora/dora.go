// Package dora TODO: package docs
package dora

import (
	"errors"
	"fmt"
	"strconv"

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
	str := string(bytes)
	return NewFromString(str)
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

func scanQueryTokens(query []rune) ([]queryToken, error) {
	var qt []queryToken
	queryLen := len(query)

	// Start at 1 to ignore the `$`, which has already been validated at this point
	for i := 1; i < queryLen-1; i++ {
		switch query[i] {
		case '.':
			// Step into the key, ex: If we were at the `.` in `.name` now we will be at `n`
			i++
			// Retrieve the selector, jump (how far to increase `i`), and err
			s, jump, _, err := parseSelector(query[i:])
			if err != nil {
				return []queryToken{}, nil
			}
			qt = append(qt, queryToken{
				accessType: "object",
				keyReq:     string(s),
			})
			i += jump - 1
		case '[':
			// Step into the key or index, ex: If we were at the `[` in `["name"]` now we are at `"`
			i++
			// Retrieve the selector, jump (how far to increase `i`), and err
			s, jump, isIndex, err := parseSelector(query[i:])
			if err != nil {
				return []queryToken{}, nil
			}
			if isIndex {
				index, err := strconv.Atoi(string(s))
				if err != nil {
					return []queryToken{}, nil
				}
				qt = append(qt, queryToken{
					accessType: "array",
					indexReq:   index,
				})
			} else {
				qt = append(qt, queryToken{
					accessType: "object",
					keyReq:     string(s),
				})
			}

			i += jump
		default:
			// TODO error?
		}
	}

	return qt, nil
}

func parseSelector(queryChunk []rune) ([]rune, int, bool, error) {
	var jump int
	var isIndex bool
	queryLen := len(queryChunk)

	if isLetter(queryChunk[jump]) {
		for isLetter(queryChunk[jump]) && jump < queryLen-1 {
			jump++
		}
		if jump < queryLen-1 {
			jump++
		}

		switch queryChunk[jump] {
		case '"':
			jump += 2 // we're finishing an object bracket notation selection "]
		case '.', '[': // we're on to another selector
			return queryChunk[0:jump], jump, isIndex, nil
		default:
			if jump == queryLen-1 { // we are on the last byte
				return queryChunk[0 : jump+1], jump, isIndex, nil
			}
			return nil, 0, isIndex, errors.New("TODO: helpful error message")
		}
	} else if isNumber(queryChunk[jump]) {
		isIndex = true
		for isNumber(queryChunk[jump]) && jump < queryLen-1 {
			jump++
		}
		return queryChunk[0:jump], jump, isIndex, nil
	}

	return nil, 0, isIndex, errors.New("TODO: helpful error message")
}

func isLetter(char rune) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func isNumber(char rune) bool {
	return '0' <= char && char <= '9'
}

func (c *Client) executeQuery() error {
	fmt.Println(c.query)
	// parseQuery into parsedQuery type deal?
	// iterate through c.program to fetch user request
	return nil
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

// queryToken represents a single "step" in each query.
// Queries are parsed into a []queryTokens to be used for exploring the JSON.
type queryToken struct {
	accessType string // object or array
	keyReq     string // a key like "name"
	indexReq   int    // an index selection like 0, 1, 2
}

/*
	JSON Object:
	{
		"name": "bradford",
		"array": ["some", "values"]
		"obj": {
			"innerKey": {
				"innerKey2": "innerValue"
			}
		}
	}

	- Query must start with $.
	$.name 							== "bradford"
	$.array 						== ["array", "values"]
	$.array[0] 						== "some"
	$.array[1] 						== "values"
	$.array[2] 						== error

	$.obj.innerKey.innerKey2 		== "innerValue"
	$.obj["innerKey"].innerKey2 	== "innerValue"
	$.obj.innerKey["innerKey2"] 	== "innerValue"
	$.obj["innerKey"]["innerKey2"] 	== "innerValue"

	---

	JSON Array:
	[
		"some",
		"values",
		{ "objKey": "objValue" }
	]

	- Query must start with $.
	$[0] 			== "some"
	$[1] 			== "values"
	$[2] 			== { "objKey": "objValue" }
	$[2].objKey 	== "objValue"
	$[2]["objKey"] 	== "objValue"

*/
