// Package parser TODO
package parser

import (
	"errors"
	"fmt"

	"github.com/bradford-hamilton/parsejson/pkg/ast"
	"github.com/bradford-hamilton/parsejson/pkg/lexer"
	"github.com/bradford-hamilton/parsejson/pkg/token"
)

// Parser holds a Lexer, its errors, the currentToken, peekToken (next token)
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

// ParseProgram parses tokens and creates an AST. It returns the RootNode
// which holds a slice of Values (and in turn, the rest of the tree)
func (p *Parser) ParseProgram() (*ast.RootNode, error) {
	var rootNode ast.RootNode

	if p.currentTokenTypeIs(token.LeftBrace) {
		for !p.currentTokenTypeIs(token.EOF) {
			val := p.parseValue()
			if val != nil {
				rootNode.Object = val.(*ast.Object)
			}
			p.nextToken()
		}

		return &rootNode, nil
	} else if p.currentTokenTypeIs(token.LeftBracket) {
		rootNode.Type = ast.ArrayRoot

		for !p.currentTokenTypeIs(token.EOF) {
			val := p.parseValue()
			if val != nil {
				rootNode.Array = val.(*ast.Array)
			}
			p.nextToken()
		}

		return &rootNode, nil
	}

	// So here do I add to the if block above and set the rootNode.Array?
	// Than I would otherwise set rootNode.Object?
	// Then have two starting points parseObject and parseArray?
	// Would parseObject just use parseArray within it?
	// Something doesn't feel quite right... need a break

	return nil, errors.New("error")
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

func (p *Parser) parseValue() ast.Value {
	switch p.currentToken.Type {
	case token.LeftBrace:
		return p.parseJSONObject()
	// case token.LeftBracket:
	// 	return p.parseJSONArray()
	// case token.String:
	// 	return p.parseJSONString()
	// case token.Minus:
	// 	return p.parseJSONumber()
	// case token.Number:
	// 	return p.parseJSONumber()
	// case token.Illegal:
	// 	return p.illegalToken()
	default:
		return nil
	}
}

func (p *Parser) parseJSONObject() ast.Value {
	obj := ast.Object{
		Type:     "Object",
		Children: []ast.Property{},
	}

	objState := ast.ObjStart

	// TODO: could be wrong, may need a subset of the tokens? We'll see
	for p.currentTokenTypeIs(token.EOF) {
		switch objState {
		case ast.ObjStart:
			if p.expectPeekType(token.LeftBrace) {
				objState = ast.ObjOpen
			} else {
				return nil
			}
		case ast.ObjOpen:
			if p.peekTokenTypeIs(token.RightBrace) {
				return obj
			}
			prop := p.parseProperty()
			obj.Children = append(obj.Children, prop)
			objState = ast.ObjProperty
		}
	}

	return obj
}

func (p *Parser) parseProperty() ast.Property {
	prop := ast.Property{
		Type: "Property",
	}

	propertyState := ast.PropertyStart

	// TODO: could be wrong, may need a subset of the tokens? We'll see
	for p.currentTokenTypeIs(token.EOF) {
		switch propertyState {
		case ast.PropertyStart:
			if p.expectPeekType(token.String) {
				key := ast.Identifier{
					Type:  "Identifier",
					Value: p.parseString(),
					Raw:   p.currentToken.Literal,
				}
				prop.Key = key
				propertyState = ast.PropertyKey
			}
		case ast.PropertyKey:
			if p.expectPeekType(token.Colon) {
				propertyState = ast.PropertyColon
			} else {
				p.errors = append(p.errors, "TODO: error here")
			}
		case ast.PropertyColon:
			val := p.parseValue()
			prop.Value = val
		}
	}

	return prop
}

func (p *Parser) parseString() ast.Value {
	res := ""

	for !p.currentTokenTypeIs(token.String) {

	}

}

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

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf(
		"Line: %d: Expected next token to be %s, got: %s instead", p.currentToken.Line, t, p.peekToken.Type,
	)
	p.errors = append(p.errors, msg)
}
