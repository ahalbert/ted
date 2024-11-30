package parser

import (
	"strconv"

	"github.com/ahalbert/fsaed/fsaed/ast"
	"github.com/ahalbert/fsaed/fsaed/lexer"
	"github.com/ahalbert/fsaed/fsaed/token"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken        token.Token
	peekToken       token.Token
	anonymousStates int
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:               l,
		errors:          []string{},
		anonymousStates: 0,
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) ParseFSA() ast.FSA {
	program := ast.FSA{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	statement := &ast.StateStatement{}
	if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.COLON) {
		statement.StateName = p.curToken.Literal
		p.nextToken()
		p.nextToken()
		statement.Action = p.parseAction()
	} else {
		statement.StateName = strconv.Itoa(p.anonymousStates)
		p.anonymousStates++
		statement.Action = p.parseAction()
	}
	return statement
}

func (p *Parser) parseAction() ast.Action {
	if p.curTokenIs(token.REGEX) {
		action := &ast.TransitionAction{Token: p.curToken, Rule: p.curToken.Literal}
		action.Rule = p.curToken.Literal
		if p.peekTokenIs(token.GOTO) {
			p.nextToken()
			p.nextToken()
			action.Target = p.curToken.Literal
		} else {
			action.Target = strconv.Itoa(p.anonymousStates)
		}
		return action
	} else if p.curTokenIs(token.DO) {
		sedAction := &ast.SedAction{Token: p.curToken, Command: p.curToken.Literal}
		p.nextToken()
		action := &ast.TransitionAction{Rule: ".*", Action: sedAction}
		if p.curTokenIs(token.GOTO) {
			p.nextToken()
			action.Target = p.curToken.Literal
		} else {
			action.Target = strconv.Itoa(p.anonymousStates)
		}
		return action
	} else {
		return nil
	}
}
