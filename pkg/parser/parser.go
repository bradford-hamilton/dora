package parser

import (
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

// nextToken sets our current token to the peek token and the peek token to
// p.lexer.NextToken() which ends up scanning and returning the next token
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}
