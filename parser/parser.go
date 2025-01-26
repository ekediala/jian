package parser

import (
	"fmt"

	"strconv"

	"github.com/ekediala/jian/ast"
	"github.com/ekediala/jian/lexer"
	"github.com/ekediala/jian/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

type (
	prefixParsefn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer          *lexer.Lexer
	curToken       token.Token
	peekToken      token.Token
	errors         []string
	prefixParsefns map[token.TokenType]prefixParsefn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := Parser{
		lexer:          l,
		errors:         make([]string, 0, 10),
		prefixParsefns: map[token.TokenType]prefixParsefn{},
		infixParseFns:  map[token.TokenType]infixParseFn{},
	}
	// initialise curToken and peekToken
	p.nextToken() // sets curToken to nil value and peekToken to first token from lexer
	p.nextToken() // sets curToken to first token from lexer and peekToken to second token from lexer

	// prefix parsing functions
	p.registerPrefixFn(token.IDENT, p.parseIdentifier)
	p.registerPrefixFn(token.INT, p.parseIntegerLiteral)
	p.registerPrefixFn(token.STRING, p.parseStringLiteral)
	p.registerPrefixFn(token.BANG, p.parsePrefixExpression)
	p.registerPrefixFn(token.MINUS, p.parsePrefixExpression)
	p.registerPrefixFn(token.TRUE, p.parseBooleanExpression)
	p.registerPrefixFn(token.FALSE, p.parseBooleanExpression)
	p.registerPrefixFn(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefixFn(token.IF, p.parseIfExpression)
	p.registerPrefixFn(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefixFn(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefixFn(token.LBRACE, p.parseHashLiteral)

	// infix parsing functions
	p.registerInfixFn(token.PLUS, p.parseInfixExpression)
	p.registerInfixFn(token.MINUS, p.parseInfixExpression)
	p.registerInfixFn(token.SLASH, p.parseInfixExpression)
	p.registerInfixFn(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixFn(token.EQ, p.parseInfixExpression)
	p.registerInfixFn(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfixFn(token.GT, p.parseInfixExpression)
	p.registerInfixFn(token.LT, p.parseInfixExpression)
	p.registerInfixFn(token.LPAREN, p.parseCallExpression)
	p.registerInfixFn(token.LBRACKET, p.parseIndexExpression)
	return &p
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := ast.HashLiteral{
		Token: p.curToken,
		Pairs: map[ast.Expression]ast.Expression{},
	}

	p.nextToken() // advance from { token

	if p.curTokenIs(token.RBRACE) {
		return &hash
	}

	key := p.parseExpression(LOWEST)
	if !p.expectPeek(token.COLON) {
		return nil
	}

	p.nextToken() // advance past colon

	value := p.parseExpression(LOWEST)
	hash.Pairs[key] = value

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // advance to comma
		p.nextToken() // advance past comma

		key = p.parseExpression(LOWEST)
		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken() // advance past colon

		value = p.parseExpression(LOWEST)
		hash.Pairs[key] = value
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return &hash
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	indexExp := &ast.IndexExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()

	indexExp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return indexExp
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{
		Token:    p.curToken,
		Elements: p.parseExpressionList(token.RBRACKET),
	}
	return array
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	// we are sitting on a left parenthesis, advance to expression
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIntegerLiteral() ast.Expression {

	v, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	return &ast.IntegerLiteral{Token: p.curToken, Value: v}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token: p.curToken,
	}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := ast.BlockStatement{
		Token: p.curToken,
	}
	block.Statements = make([]ast.Statement, 0, 10)

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return &block
}

func (p *Parser) parseIfExpression() ast.Expression {
	exp := ast.IfExpression{
		Token: p.curToken,
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	exp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		exp.Alternative = p.parseBlockStatement()
	}

	return &exp
}

func (p *Parser) parseBooleanExpression() ast.Expression {
	b := ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
	return &b
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return &expression
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	// get next token and call parse expression on it
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return &expression
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	parameters := make([]*ast.Identifier, 0, 5)

	if p.peekTokenIs(token.RPAREN) {
		// function has no parameters
		p.nextToken()
		return parameters
	}

	p.nextToken()

	parameter := ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	parameters = append(parameters, &parameter)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // go to comma
		p.nextToken() // go to next parameter
		parameter := ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}

		parameters = append(parameters, &parameter)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return parameters
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	exp := ast.FunctionLiteral{
		Token: p.curToken,
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	exp.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Body = p.parseBlockStatement()

	return &exp
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	exp := ast.CallExpression{Token: p.curToken, Function: fn}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return &exp
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixFn, ok := p.prefixParsefns[p.curToken.Type]
	if !ok {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefixFn()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infixFn := p.infixParseFns[p.peekToken.Type]
		if infixFn == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infixFn(leftExp)
	}

	return leftExp
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return &stmt
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := ast.NewProgram(10)
	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek checks if the next token (peekToken) matches the expected token type.
// If it matches, it advances the parser to the next token and returns true.
// If it doesn't match, it returns false and adds an error to the parser
// without advancing the parser.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) registerInfixFn(t token.TokenType, fn infixParseFn) {
	p.infixParseFns[t] = fn
}

func (p *Parser) registerPrefixFn(t token.TokenType, fn prefixParsefn) {
	p.prefixParsefns[t] = fn
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	var list = []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // get us to comma
		p.nextToken() // get us to next expression
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}
