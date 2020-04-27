package dora

import (
	"fmt"
	"strconv"

	"github.com/bradford-hamilton/dora/pkg/danger"
)

// The available accessTypes for a dora query
const (
	ObjectAccess accessType = iota
	ArrayAccess
)

type accessType int

// queryToken represents a single "step" in each query.
// Queries are parsed into a []queryTokens to be used for exploring the JSON.
type queryToken struct {
	accessType accessType // ObjectAccess or ArrayAccess
	key        string     // a key like "name"
	index      int        // an index selection like 0, 1, 2
}

// scanQueryTokens scans a users query input into a collection of queryTokens.
// Dora's query syntax is very straight forward, here is a quick BNF-like representation:
//    <dora-query>  ::= <querystring>
//    <querystring> ::= "<query>,*"
//    <query>       ::= "[<int>]" | "." + <string>
func scanQueryTokens(query []byte) ([]queryToken, error) {
	var qts []queryToken
	queryLen := len(query)

	// Start at 1 to ignore the `$`, which has already been validated at this point.
	for i := 1; i < queryLen-1; i++ {
		switch query[i] {
		case '.':
			// Step into the key, ex: - If we were at the `.` in `.name` this bumps us to `n`.
			i++

			// Retrieve the selector and how far to increase `i` (jump).
			s, jump, _, err := parseObjSelector(query[i:])
			if err != nil {
				return []queryToken{}, err
			}

			// Append our new query token and adjust the jump.
			qts = append(qts, queryToken{accessType: ObjectAccess, key: danger.BytesToString(s)})
			i += jump - 1
		case '[':
			// Step into the index, ex: - If we were at the `[` in `[123]` this bumps us to `1`
			i++

			// Retrieve the selector and how far to increase `i` (jump).
			s, jump, err := parseArraySelector(query[i:])
			if err != nil {
				return []queryToken{}, err
			}

			// The array selector is an int, and we want assert that.
			index, err := strconv.Atoi(danger.BytesToString(s))
			if err != nil {
				return []queryToken{}, err
			}

			// Append our new query token and adjust the jump
			qts = append(qts, queryToken{accessType: ArrayAccess, index: index})
			i += jump
		default:
			return []queryToken{}, errSelectorSytax(string(query[i]))
		}
	}

	return qts, nil
}

// parseObjSelector consumes the property key, sets the `jump` index to right after it, and returns the sliced chunk.
func parseObjSelector(queryChunk []byte) ([]byte, int, bool, error) {
	var jump int
	var isIndex bool
	queryLen := len(queryChunk)

	if isPropertyKey(queryChunk[jump]) {
		// Consume the property key name and increment the jump to point to after the property key
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

// parseArraySelector consumes the array index request, sets the `jump` index to right after it, and returns the sliced chunk.
func parseArraySelector(queryChunk []byte) ([]byte, int, error) {
	var jump int
	queryLen := len(queryChunk)

	if isNumber(queryChunk[jump]) {
		// Consume the index and return it along with the jump
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

func isPropertyKey(char byte) bool {
	return isLetter(char) || isNumber(char)
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func isNumber(char byte) bool {
	return '0' <= char && char <= '9'
}
