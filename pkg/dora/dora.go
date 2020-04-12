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

// queryToken represents a single "step" in each query.
// Queries are parsed into a []queryTokens to be used for exploring the JSON.
type queryToken struct {
	accessType string // object or array
	keyReq     string // a key like "name"
	indexReq   int    // an index selection like 0, 1, 2
}

func (c *Client) parseQuery() error {
	tokens, err := scanQueryTokens(c.query)
	if err != nil {
		return err
	}
	c.parsedQuery = tokens
	return nil
}

// scanQueryTokens scans a users query input into a collection of queryTokens.
// Dora's query syntax is very straight forward, here is a quick BNF-like representation:
//	  <dora-query>  ::= <querystring>
//    <querystring> ::= "<query>,*"
//    <query>       ::= "[<int>]" | ".<string>"
func scanQueryTokens(query []rune) ([]queryToken, error) {
	var qts []queryToken
	queryLen := len(query)

	// Start at 1 to ignore the `$`, which has already been validated at this point.
	for i := 1; i < queryLen-1; i++ {
		switch query[i] {
		case '.':
			// Step into the key, ex:
			// - If we were at the `.` in `.name` this bumps us to `n`.
			i++

			// Retrieve the selector and how far to increase `i` (jump).
			s, jump, _, err := parseObjSelector(query[i:])
			if err != nil {
				return []queryToken{}, err
			}

			// Append our new query token and adjust the jump.
			qts = append(qts, queryToken{accessType: "object", keyReq: string(s)})
			i += jump - 1
		case '[':
			// Step into the index, ex:
			// - If we were at the `[` in `[123]` this bumps us to `1`
			i++

			// Retrieve the selector and how far to increase `i` (jump).
			s, jump, err := parseArraySelector(query[i:])
			if err != nil {
				return []queryToken{}, err
			}

			// The array selector is an int, and we want assert that.
			index, err := strconv.Atoi(string(s))
			if err != nil {
				return []queryToken{}, err
			}

			// Append our new query token and adjust the jump
			qts = append(qts, queryToken{accessType: "array", indexReq: index})
			i += jump
		default:
			return []queryToken{}, errSelectorSytax(string(query[i]))
		}
	}

	return qts, nil
}

func parseObjSelector(queryChunk []rune) ([]rune, int, bool, error) {
	var jump int
	var isIndex bool
	queryLen := len(queryChunk)

	// Consume the key name
	if isPropertyKey(queryChunk[jump]) {
		for isPropertyKey(queryChunk[jump]) && jump < queryLen-1 {
			jump++
		}

		// After consuming the property key above, we should be onto the next selector.
		// This has to either be another object or an array
		if queryChunk[jump] == '.' || queryChunk[jump] == '[' {
			return queryChunk[0:jump], jump, isIndex, nil
		} else if jump == queryLen-1 {
			// we are on the last byte
			return queryChunk[0 : jump+1], jump, isIndex, nil
		}

		return nil, 0, isIndex, errSelectorSytax(string(queryChunk[jump]))
	}

	return nil, 0, isIndex, fmt.Errorf(
		"Error parsing object selector within query. Expected string, but started with %s",
		string(queryChunk[jump]),
	)
}

func parseArraySelector(queryChunk []rune) ([]rune, int, error) {
	var jump int
	queryLen := len(queryChunk)

	// Consume the index and return it along with the jump
	if isNumber(queryChunk[jump]) {
		for isNumber(queryChunk[jump]) && jump < queryLen-1 {
			jump++
		}
		return queryChunk[0:jump], jump, nil
	}

	return nil, 0, fmt.Errorf(
		"Error parsing array selector within query. Expected an int, but started with %s",
		string(queryChunk[jump]),
	)
}

func (c *Client) executeQuery() error {
	// TODO: actually execute the query next
	fmt.Println(c.query)
	// parseQuery into parsedQuery type deal?
	// iterate through c.program to fetch user request
	return nil
}

func isPropertyKey(char rune) bool {
	return isLetter(char) || isNumber(char)
}

func isLetter(char rune) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func isNumber(char rune) bool {
	return '0' <= char && char <= '9'
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
