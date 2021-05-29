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
	p := Parser{lexer: l}

	// Read two tokens, so currentToken and peekToken are both set.
	p.nextToken()
	p.nextToken()

	return &p
}

// ParseJSON parses tokens and creates an AST. It returns the RootNode
// which holds a slice of Values (and in turn, the rest of the tree)
func (p *Parser) ParseJSON() (ast.RootNode, error) {
	var rootNode ast.RootNode
	if p.currentTokenTypeIs(token.LeftBracket) {
		rootNode.Type = ast.ArrayRoot
	}

	val := p.parseValue()
	if val.Content == nil {
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
	value := ast.Value{
		PrefixStructure: p.parseStructure(),
	}

	switch p.currentToken.Type {
	case token.LeftBrace:
		value.Content = p.parseJSONObject()
	case token.LeftBracket:
		value.Content = p.parseJSONArray()
	default:
		value.Content = p.parseJSONLiteral()
	}

	value.SuffixStructure = p.parseStructure()

	return value
}

func (p *Parser) parseArrayItem() ast.ArrayItem {
	arrayItem := ast.ArrayItem{
		Type:            ast.ArrayItemType,
		PrefixStructure: p.parseStructure(),
	}

	switch p.currentToken.Type {
	case token.LeftBrace:
		arrayItem.Value = p.parseJSONObject()
	case token.LeftBracket:
		arrayItem.Value = p.parseJSONArray()
	default:
		arrayItem.Value = p.parseJSONLiteral()
	}

	arrayItem.PostValueStructure = p.parseStructure()

	return arrayItem

}

// parseJSONObject is called when an open left brace `{` token is found
func (p *Parser) parseJSONObject() ast.ValueContent {
	obj := ast.Object{Type: ast.ObjectType}
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
				obj.Children[len(obj.Children)-1].HasCommaSeparator = true
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
			structure := p.parseStructure()
			if p.currentTokenTypeIs(token.RightBrace) {
				obj.SuffixStructure = structure
				p.nextToken()
				obj.End = p.currentToken.End
				return obj
			}
			prop := p.parseProperty()
			prop.PrefixStructure = append(structure, prop.PrefixStructure...)
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
func (p *Parser) parseJSONArray() ast.ValueContent {
	array := ast.Array{Type: ast.ArrayType}
	arrayState := ast.ArrayStart
	array.PrefixStructure = p.parseStructure()

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
			arrayItem := p.parseArrayItem()
			array.Children = append(array.Children, arrayItem)
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
				array.Children[len(array.Children)-1].HasCommaSeparator = true
				arrayState = ast.ArrayComma
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf(
					"Error parsing property. Expected RightBrace or Comma token, got: %s",
					p.currentToken.Literal,
				))
			}
		case ast.ArrayComma:
			structure := p.parseStructure()
			if p.currentTokenTypeIs(token.RightBracket) {
				array.SuffixStructure = structure
				p.nextToken()
				array.End = p.currentToken.End
				return array
			}
			arrayItem := p.parseArrayItem()
			arrayItem.PrefixStructure = append(structure, array.PrefixStructure...)
			array.Children = append(array.Children, arrayItem)
			arrayState = ast.ArrayValue
		}
	}
	array.End = p.currentToken.Start
	array.SuffixStructure = p.parseStructure()
	return array
}

// parseJSONLiteral switches on the current token's type, sets the Value on a return val and returns it.
func (p *Parser) parseJSONLiteral() ast.Literal {
	val := ast.Literal{Type: ast.LiteralType}

	// Regardless of what the current token type is - after it's been assigned, we must consume the token
	defer p.nextToken()

	switch p.currentToken.Type {
	case token.String:
		val.ValueType = ast.StringLiteralValueType
		val.Delimiter = p.currentToken.Prefix
		val.Value = p.parseString()
		return val
	case token.Number:
		val.ValueType = ast.NumberLiteralValueType
		ct := p.currentToken.Literal
		val.OriginalRendering = ct

		// Attempt to parse as an integer first, then float
		i, err := strconv.ParseInt(ct, 10, 64)
		if err == nil {
			val.Value = i
			return val
		}
		f, err := strconv.ParseFloat(ct, 64)
		if err != nil {
			p.parseError("error parsing JSON number, incorrect syntax")
			val.Value = ct
			return val
		}
		val.Value = f

		return val
	case token.True:
		val.ValueType = ast.BooleanLiteralValueType
		val.Value = true
		return val
	case token.False:
		val.ValueType = ast.BooleanLiteralValueType
		val.Value = false
		return val
	default:
		val.ValueType = ast.NullLiteralValueType
		val.Value = "null"
		return val
	}
}

// parseProperty is used to parse an object property and in doing so handles setting the `key`:`value` pair.
func (p *Parser) parseProperty() ast.Property {
	prop := ast.Property{Type: ast.PropertyType}
	propertyState := ast.PropertyStart

	for !p.currentTokenTypeIs(token.EOF) {
		switch propertyState {
		case ast.PropertyStart:
			prop.PrefixStructure = p.parseStructure()
			if p.currentTokenTypeIs(token.String) {
				key := ast.Identifier{
					Type:      ast.IdentifierType,
					Value:     p.parseString(),
					Delimiter: p.currentToken.Prefix,
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
			prop.PostKeyStructure = p.parseStructure()

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
			prop.PreValueStructure = p.parseStructure()
			val := p.parseValue()
			prop.Value = val
			propertyState = ast.PropertyValue
		case ast.PropertyValue:
			prop.PostValueStructure = p.parseStructure()
			return prop
		}
	}

	return prop
}

func (p *Parser) parseStructure() []ast.StructuralItem {
	var result []ast.StructuralItem
	for {
		switch p.currentToken.Type {
		case token.Whitespace, token.BlockComment, token.LineComment:
			value := p.currentToken.Prefix + p.currentToken.Literal + p.currentToken.Suffix
			result = append(result, ast.StructuralItem{Value: value})
			p.nextToken()
		default:
			return result
		}
	}
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
