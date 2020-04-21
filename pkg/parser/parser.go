// Package parser TODO
package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bradford-hamilton/dora/pkg/ast"
	"github.com/bradford-hamilton/dora/pkg/lexer"
	"github.com/bradford-hamilton/dora/pkg/token"
)

// Parser holds a Lexer, errors, the currentToken, and the peek peekToken (next token).
// Parser methods handle iterating through tokens and building and AST.
type Parser struct {
	lexer        *lexer.Lexer
	errors       []string
	currentToken token.Token
	peekToken    token.Token
}

// New takes a Lexer, creates a Parser with that Lexer, sets the current and
// peek tokens, and returns the Parser.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}

	// Read two tokens, so currentToken and peekToken are both set.
	p.nextToken()
	p.nextToken()

	return p
}

// ParseJSON parses tokens and creates an AST. It returns the RootNode
// which holds a slice of Values (and in turn, the rest of the tree)
func (p *Parser) ParseJSON() (ast.RootNode, error) {
	var rootNode ast.RootNode
	if p.currentTokenTypeIs(token.LeftBracket) {
		rootNode.Type = ast.ArrayRoot
	}

	val := p.parseValue()
	if val == nil {
		p.parseError(fmt.Sprintf(
			"Error parsing JSON expected a value, got: %v:",
			p.currentToken.Literal,
		))
		return ast.RootNode{}, errors.New(p.Errors())
	}
	rootNode.RootValue = &val

	return rootNode, nil
}

// nextToken sets our current token to the peek token and the peek token to
// p.lexer.NextToken() which ends up scanning and returning the next token
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) currentTokenTypeIs(t token.Type) bool {
	return p.currentToken.Type == t
}

// parseValue is our dynamic entrypoint to parsing JSON values. All scenarios for
// this parser fall under these 3 actions.
func (p *Parser) parseValue() ast.Value {
	switch p.currentToken.Type {
	case token.LeftBrace:
		return p.parseJSONObject()
	case token.LeftBracket:
		return p.parseJSONArray()
	default:
		return p.parseJSONLiteral()
	}
}

// parseJSONObject is called when an open left brace `{` token is found
func (p *Parser) parseJSONObject() ast.Value {
	obj := ast.Object{Type: "Object"}
	objState := ast.ObjStart

	for !p.currentTokenTypeIs(token.EOF) {
		switch objState {
		case ast.ObjStart:
			if p.currentTokenTypeIs(token.LeftBrace) {
				objState = ast.ObjOpen
				obj.Start = p.currentToken.Start
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf(
					"Error parsing JSON object Expected `{` token, got: %s",
					p.currentToken.Literal,
				))
				return nil
			}
		case ast.ObjOpen:
			if p.currentTokenTypeIs(token.RightBrace) {
				p.nextToken()
				obj.End = p.currentToken.End
				return obj
			}
			prop := p.parseProperty()
			obj.Children = append(obj.Children, prop)
			objState = ast.ObjProperty
		case ast.ObjProperty:
			if p.currentTokenTypeIs(token.RightBrace) {
				p.nextToken()
				obj.End = p.currentToken.Start
				return obj
			} else if p.currentTokenTypeIs(token.Comma) {
				objState = ast.ObjComma
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf(
					"Error parsing property. Expected RightBrace or Comma token, got: %s",
					p.currentToken.Literal,
				))
				return nil
			}
		case ast.ObjComma:
			prop := p.parseProperty()
			if prop.Value != nil {
				obj.Children = append(obj.Children, prop)
				objState = ast.ObjProperty
			}
		}
	}

	obj.End = p.currentToken.Start

	return obj
}

// parseJSONArray is called when an open left bracket `[` token is found
func (p *Parser) parseJSONArray() ast.Value {
	array := ast.Array{Type: "Array"}
	arrayState := ast.ArrayStart

	for !p.currentTokenTypeIs(token.EOF) {
		switch arrayState {
		case ast.ArrayStart:
			if p.currentTokenTypeIs(token.LeftBracket) {
				array.Start = p.currentToken.Start
				arrayState = ast.ArrayOpen
				p.nextToken()
			}
		case ast.ArrayOpen:
			if p.currentTokenTypeIs(token.RightBracket) {
				array.End = p.currentToken.End
				p.nextToken()
				return array
			}
			val := p.parseValue()
			array.Children = append(array.Children, val)
			arrayState = ast.ArrayValue
			if p.peekTokenTypeIs(token.RightBracket) {
				p.nextToken()
			}
		case ast.ArrayValue:
			if p.currentTokenTypeIs(token.RightBracket) {
				array.End = p.currentToken.End
				p.nextToken()
				return array
			} else if p.currentTokenTypeIs(token.Comma) {
				arrayState = ast.ArrayComma
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf(
					"Error parsing property. Expected RightBrace or Comma token, got: %s",
					p.currentToken.Literal,
				))
			}
		case ast.ArrayComma:
			val := p.parseValue()
			array.Children = append(array.Children, val)
			arrayState = ast.ArrayValue
		}
	}
	array.End = p.currentToken.Start
	return array
}

// parseJSONLiteral switches on the current token's type, sets the Value on a return val and returns it.
func (p *Parser) parseJSONLiteral() ast.Literal {
	val := ast.Literal{Type: "Literal"}

	// Regardless of what the current token type is - after it's been assigned, we must consume the token
	defer p.nextToken()

	switch p.currentToken.Type {
	case token.String:
		val.Value = p.parseString()
		return val
	case token.Number:
		v, _ := strconv.Atoi(p.currentToken.Literal)
		val.Value = v
		return val
	case token.True:
		val.Value = true
		return val
	case token.False:
		val.Value = false
		return val
	default:
		val.Value = "null"
		return val
	}
}

// parseProperty is used to parse an object property and in doing so handles setting the `key`:`value` pair.
func (p *Parser) parseProperty() ast.Property {
	prop := ast.Property{Type: "Property"}
	propertyState := ast.PropertyStart

	for !p.currentTokenTypeIs(token.EOF) {
		switch propertyState {
		case ast.PropertyStart:
			if p.currentTokenTypeIs(token.String) {
				key := ast.Identifier{
					Type:  "Identifier",
					Value: p.parseString(),
				}
				prop.Key = key
				propertyState = ast.PropertyKey
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf(
					"Error parsing property start. Expected String token, got: %s",
					p.currentToken.Literal,
				))
			}
		case ast.PropertyKey:
			if p.currentTokenTypeIs(token.Colon) {
				propertyState = ast.PropertyColon
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf(
					"Error parsing property. Expected Colon token, got: %s",
					p.currentToken.Literal,
				))
			}
		case ast.PropertyColon:
			val := p.parseValue()
			prop.Value = val
			return prop
		}
	}

	return prop
}

// TODO: all the tedius ecaping, etc still needs to be applied here
func (p *Parser) parseString() string {
	return p.currentToken.Literal
}

// expectPeekType checks the next token type against the one passed in. If it matches,
// we call p.nextToken() to set us to the expected token and return true. If the expected
// type does not match, we add a peek error and return false.
func (p *Parser) expectPeekType(t token.Type) bool {
	if p.peekTokenTypeIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekTokenTypeIs(t token.Type) bool {
	return p.peekToken.Type == t
}

// peekError is a small wrapper to add a peek error to our parser's errors field.
func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("Line: %d: Expected next token to be %s, got: %s instead", p.currentToken.Line, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// parseError is very similar to `peekError`, except it simply takes a string message that
// gets appended to the parser's errors
func (p *Parser) parseError(msg string) {
	p.errors = append(p.errors, msg)
}

// Errors is simply a helper function that returns the parser's errors
func (p *Parser) Errors() string {
	return strings.Join(p.errors, ", ")
}
