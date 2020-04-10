// Package parser TODO
package parser

import (
	"errors"
	"fmt"
	"strconv"

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
		val := p.parseJSONObject()
		if val != nil {
			rootNode.Object = &val
		}

		return &rootNode, nil
	} else if p.currentTokenTypeIs(token.LeftBracket) {
		rootNode.Type = ast.ArrayRoot
		val := p.parseJSONArray()
		if val != nil {
			rootNode.Array = &val
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
	case token.LeftBracket:
		return p.parseJSONArray()
	default:
		return p.parseJSONLiteral()
	}
}

func (p *Parser) parseJSONObject() ast.Value {
	obj := ast.Object{Type: "Object"}
	objState := ast.ObjStart

	// TODO: could be wrong, may need a subset of the tokens? We'll see
	for !p.currentTokenTypeIs(token.EOF) {
		switch objState {
		case ast.ObjStart:
			if p.currentTokenTypeIs(token.LeftBrace) {
				objState = ast.ObjOpen
				p.nextToken()
			} else {
				// add to errors
				return nil
			}
		case ast.ObjOpen:
			if p.peekTokenTypeIs(token.RightBrace) {
				p.nextToken()
				return obj
			}
			prop := p.parseProperty()
			obj.Children = append(obj.Children, prop)
			objState = ast.ObjProperty

			// if !p.currentTokenTypeIs(token.Comma) {
			// 	p.nextToken()
			// }

		case ast.ObjProperty:
			if p.currentTokenTypeIs(token.RightBrace) {
				p.nextToken()
				return obj
			} else if p.currentTokenTypeIs(token.Comma) {
				objState = ast.ObjComma
				p.nextToken()
			} else {
				// error
				fmt.Println("err")
			}
		case ast.ObjComma:
			prop := p.parseProperty()
			if prop.Value != nil {
				obj.Children = append(obj.Children, prop)
				objState = ast.ObjProperty
			}
			// p.nextToken()
		}
	}

	return obj
}

func (p *Parser) parseJSONArray() ast.Value {
	array := ast.Array{Type: "Array"}
	arrayState := ast.ArrayStart

	// TODO: could be wrong, may need a subset of the tokens? We'll see
	for !p.currentTokenTypeIs(token.EOF) {
		switch arrayState {
		case ast.ArrayStart:
			if p.currentTokenTypeIs(token.LeftBracket) {
				arrayState = ast.ArrayOpen
				p.nextToken()
			}
		case ast.ArrayOpen:
			if p.currentTokenTypeIs(token.RightBracket) {
				p.nextToken()
				return array
			}
			val := p.parseValue()
			array.Children = append(array.Children, val)
			arrayState = ast.ArrayValue

			// NOTE: this is what was helping in some scenarios but not others
			// if p.peekTokenTypeIs(token.RightBracket) {
			// 	p.nextToken()
			// }

		case ast.ArrayValue:
			if p.currentTokenTypeIs(token.RightBracket) {
				p.nextToken()
				return array
			} else if p.currentTokenTypeIs(token.Comma) {
				arrayState = ast.ArrayComma
				p.nextToken()
			}
		case ast.ArrayComma:
			val := p.parseValue()
			array.Children = append(array.Children, val)
			arrayState = ast.ArrayValue
			// p.nextToken()
		}
	}

	return array
}

func (p *Parser) parseJSONLiteral() ast.Literal {
	val := ast.Literal{Type: "Literal"}

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
	case token.Null:
		val.Value = nil
		return val
	default:
		val.Value = nil
		return val
	}
}

func (p *Parser) parseProperty() ast.Property {
	prop := ast.Property{Type: "Property"}
	propertyState := ast.PropertyStart

	// TODO: could be wrong, may need a subset of the tokens? We'll see
	for !p.currentTokenTypeIs(token.EOF) {
		switch propertyState {
		case ast.PropertyStart:
			if p.currentTokenTypeIs(token.String) {
				key := ast.Identifier{
					Type:  "Identifier",
					Value: p.parseString(),
					Raw:   p.currentToken.Literal,
				}
				prop.Key = key
				propertyState = ast.PropertyKey
				p.nextToken()
			} else {
				// error
			}
		case ast.PropertyKey:
			if p.currentTokenTypeIs(token.Colon) {
				propertyState = ast.PropertyColon
				p.nextToken()
			} else {
				p.errors = append(p.errors, "TODO: error here")
			}
		case ast.PropertyColon:
			val := p.parseValue()
			prop.Value = val
			p.nextToken()
			return prop
		}
	}

	return prop
}

func (p *Parser) parseString() string {
	// TODO: all the tedius ecaping, etc still needs to be applied here
	return p.currentToken.Literal
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
