package dora

import (
	"fmt"
	"strconv"
)

// The availble accessTypes
const (
	ObjectAccess accessType = iota
	ArrayAccess
)

type accessType int

// queryToken represents a single "step" in each query.
// Queries are parsed into a []queryTokens to be used for exploring the JSON.
type queryToken struct {
	accessType accessType // ObjectAccess or ArrayAccess
	keyReq     string     // a key like "name"
	indexReq   int        // an index selection like 0, 1, 2
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
			qts = append(qts, queryToken{accessType: ObjectAccess, keyReq: string(s)})
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
			qts = append(qts, queryToken{accessType: ArrayAccess, indexReq: index})
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

func isPropertyKey(char rune) bool {
	return isLetter(char) || isNumber(char)
}

func isLetter(char rune) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func isNumber(char rune) bool {
	return '0' <= char && char <= '9'
}
