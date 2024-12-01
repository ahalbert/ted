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
	AnonymousStates int
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:               l,
		errors:          []string{},
		AnonymousStates: 1,
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
			name := stmt.(*ast.StateStatement).StateName
			program.Statements = append(program.Statements, stmt)
			p.nextToken()
			for p.curTokenIs(token.SEMICOLON) {
				p.nextToken()
				stmt := &ast.StateStatement{StateName: name, Action: p.parseAction()}
				program.Statements = append(program.Statements, stmt)
				p.nextToken()
			}
		} else {
			p.nextToken()
		}
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
		statement.StateName = strconv.Itoa(p.AnonymousStates)
		p.AnonymousStates++
		statement.Action = p.parseAction()
	}
	return statement
}

func (p *Parser) parseAction() ast.Action {
	var action ast.Action
	switch p.curToken.Type {
	case token.REGEX:
		action = p.parseRegexAction()
	case token.GOTO:
		action = p.parseGotoAction()
	case token.DO:
		action = p.parseDoAction()
	case token.PRINT:
		action = p.parsePrintAction()
	default:
		return nil
	}
	return action
}

func (p *Parser) parseRegexAction() ast.RegexAction {
	action := ast.RegexAction{Rule: p.curToken.Literal}
	action.Action = p.parseGotoAction()
	return action
}

func (p *Parser) parseGotoAction() ast.GotoAction {
	return ast.GotoAction{Target: p.helpParseGoto()}
}

func (p *Parser) helpParseGoto() string {
	switch p.peekToken.Type {
	case token.GOTO:
		p.nextToken()
		p.nextToken()
		if p.curTokenIs(token.IDENT) {
			return p.curToken.Literal
		} else {
			panic("parse error")
		}
	default:
		return strconv.Itoa(p.AnonymousStates)
	}
}

func (p *Parser) parseDoAction() ast.DoSedAction {
	action := ast.DoSedAction{Command: p.curToken.Literal}
	return action
}

func (p *Parser) parsePrintAction() ast.PrintAction {
	return ast.PrintAction{}
}
