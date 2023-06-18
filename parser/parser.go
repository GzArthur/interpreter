package parser

import (
	"fmt"
	"github.com/GzArthur/interpreter/ast"
	"github.com/GzArthur/interpreter/lexer"
	"github.com/GzArthur/interpreter/token"
	"strconv"
)

const (
	// Priority definition
	_           int = iota
	LOWEST          // doesn't exist yet
	EQUALS          // ==
	LESSGREATER     // > or <
	SUM             // + or -
	PRODUCT         // * or /
	PREFIX          // !x or -x
	CALL            // fn(x)
)

var precedences = map[token.Type]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(node ast.Expression) ast.Expression
)

type Parser struct {
	l              *lexer.Lexer
	currToken      token.Token
	peekToken      token.Token
	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
	errors         []string // collect exception info during parsing
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		errors:         []string{},
		prefixParseFns: make(map[token.Type]prefixParseFn),
		infixParseFns:  make(map[token.Type]infixParseFn),
	}
	// initialize the currToken and the peekToken
	p.nextToken()
	p.nextToken()
	// register prefix functions
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	// register infix functions
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallFunction)
	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}
	for !p.expectCurrTokenType(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken}

	if !p.expectPeekIs(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeekIs(token.ASSIGN) {
		return nil
	}

	p.nextToken() // skip ASSIGN

	stmt.Value = p.parseExpression(LOWEST)

	if p.expectPeekTokenType(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.expectPeekTokenType(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.expectPeekTokenType(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		p.CollectPrefixParseFnError(p.currToken.Type)
		return nil
	}
	leftExpr := prefix()
	for !p.expectPeekTokenType(token.SEMICOLON) && precedence < p.peekPrecedence() {
		// precedence express the current right constraint capacity
		// p.peekPrecedence() express the current left constraint capacity
		// precedence < p.peekPrecedence() condition checks if the left constraint capacity of the next operator or lexical unit
		// is stronger than the current right constraint capacity.
		// if so, the current parsed content will be fused from left to right by the next operator
		// and passed to the next operator's infixParseFn
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			break
		}
		p.nextToken()
		leftExpr = infix(leftExpr)
	}
	return leftExpr
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.currToken,
		Value: p.expectCurrTokenType(token.TRUE),
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	expr := &ast.Integer{Token: p.currToken}
	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("could not parse %q as integer", p.currToken.Literal))
		return nil
	}
	expr.Value = value
	return expr
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}
	p.nextToken()
	expr.RightExpr = p.parseExpression(PREFIX)
	return expr
}

func (p *Parser) parseInfixExpression(leftExpr ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.currToken,
		LeftExpr: leftExpr,
		Operator: p.currToken.Literal,
	}
	precedence := p.currPrecedence()
	p.nextToken()
	expr.RightExpr = p.parseExpression(precedence)
	return expr
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expr := p.parseExpression(LOWEST)

	if !p.expectPeekIs(token.RPAREN) {
		return nil
	}
	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{Token: p.currToken}

	if !p.expectPeekIs(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expr.Condition = p.parseExpression(LOWEST)

	if !p.expectPeekIs(token.RPAREN) {
		return nil
	}
	if !p.expectPeekIs(token.LBRACE) {
		return nil
	}

	expr.Consequence = p.parseBlockStatement()
	if p.expectPeekTokenType(token.ELSE) {
		p.nextToken()
		if !p.expectPeekIs(token.LBRACE) {
			return nil
		}
		expr.Alternative = p.parseBlockStatement()
	}
	return expr
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	expr := &ast.FunctionLiteral{Token: p.currToken}
	if !p.expectPeekIs(token.LPAREN) {
		return nil
	}
	expr.Parameters = p.parseFunctionParameters()
	if !p.expectPeekIs(token.LBRACE) {
		return nil
	}
	expr.Body = p.parseBlockStatement()
	return expr
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var identifiers []*ast.Identifier
	if p.expectPeekTokenType(token.RPAREN) {
		// case fn()
		p.nextToken()
		return identifiers
	}

	p.nextToken()
	for {
		identifier := &ast.Identifier{
			Token: p.currToken,
			Value: p.currToken.Literal,
		}
		identifiers = append(identifiers, identifier)
		if p.expectPeekTokenType(token.COMMA) {
			p.nextToken()
			p.nextToken()
			continue
		}
		break
	}
	if !p.expectPeekIs(token.RPAREN) {
		return nil
	}
	return identifiers
}

func (p *Parser) parseCallFunction(function ast.Expression) ast.Expression {
	expr := &ast.CallExpression{Token: p.currToken, Function: function}
	expr.Arguments = p.parseCallArguments()
	return expr
}

func (p *Parser) parseCallArguments() []ast.Expression {
	var args []ast.Expression
	if p.expectPeekTokenType(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	for {
		args = append(args, p.parseExpression(LOWEST))
		if p.expectPeekTokenType(token.COMMA) {
			p.nextToken()
			p.nextToken()
			continue
		}
		break
	}
	if !p.expectPeekIs(token.RPAREN) {
		return nil
	}
	return args
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	bStmt := &ast.BlockStatement{
		Token:      p.currToken,
		Statements: []ast.Statement{},
	}
	p.nextToken()
	for !p.expectCurrTokenType(token.RBRACE) && !p.expectCurrTokenType(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			bStmt.Statements = append(bStmt.Statements, stmt)
		}
		p.nextToken()
	}
	return bStmt
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) CollectPrefixParseFnError(t token.Type) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) CollectPeekTokenTypeError(expectedType token.Type) {
	msg := fmt.Sprintf("expected next token type to be %s, got %s instead", expectedType, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) expectPeekIs(t token.Type) bool {
	if p.expectPeekTokenType(t) {
		p.nextToken()
		return true
	} else {
		p.CollectPeekTokenTypeError(t)
		return false
	}
}

func (p *Parser) expectPeekTokenType(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectCurrTokenType(t token.Type) bool {
	return p.currToken.Type == t
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) currPrecedence() int {
	if precedence, ok := precedences[p.currToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.ReadToken()
}

func (p *Parser) registerPrefix(t token.Type, fn prefixParseFn) {
	p.prefixParseFns[t] = fn
}

func (p *Parser) registerInfix(t token.Type, fn infixParseFn) {
	p.infixParseFns[t] = fn
}
