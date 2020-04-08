// Package parser TODO
package parser

import (
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

	if p.currentTokenTypeIs(token.LeftBracket) {
		rootNode.Type = ast.ArrayRoot
	}

	// So here do I add to the if block above and set the rootNode.Array?
	// Than I would otherwise set rootNode.Object?
	// Then have two starting points parseObject and parseArray?
	// Would parseObject just use parseArray within it?
	// Something doesn't feel quite right... need a break

	for !p.currentTokenTypeIs(token.EOF) {
		// val := p.parseValue()
		// if val != nil {
		// 	rootNode.Values = append(rootNode.Values, val)
		// }
		p.nextToken()
	}

	return &rootNode, nil
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
	// switch p.currentToken.Type {
	// case token.LeftBrace:
	// 	return p.parseJSONObject()
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
	// default:
	// 	return p.parseJSONValue()
	// }

	return nil
}
