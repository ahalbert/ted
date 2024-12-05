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

	curToken                       token.Token
	peekToken                      token.Token
	AnonymousStates                int
	GOTOsThatNeedNextStateAssigned []*ast.GotoAction
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:                              l,
		errors:                         []string{},
		AnonymousStates:                1,
		GOTOsThatNeedNextStateAssigned: []*ast.GotoAction{},
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
		program.Statements = append(program.Statements, stmt)
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	statement := &ast.StateStatement{}
	if p.curTokenIs(token.LABEL) {
		statement.StateName = p.curToken.Literal
		p.nextToken()
	} else {
		statement.StateName = strconv.Itoa(p.AnonymousStates)
		p.AnonymousStates++
	}
	for _, action := range p.GOTOsThatNeedNextStateAssigned {
		(*action).Target = statement.StateName
	}
	p.GOTOsThatNeedNextStateAssigned = []*ast.GotoAction{}
	statement.Action = p.parseAction()
	return statement
}

func (p *Parser) parseAction() ast.Action {
	var action ast.Action
	switch p.curToken.Type {
	case token.LBRACE:
		action = p.parseActionBlock()
	case token.REGEX:
		action = p.parseRegexAction()
	case token.GOTO:
		action = p.parseGotoAction()
	case token.RESET:
		action = p.parseGotoAction()
	case token.DO:
		action = p.parseDoAction()
	case token.PRINT:
		action = p.parsePrintAction()
	default:
		panic("parser error: expected action")
	}
	return action
}

func (p *Parser) parseActionBlock() *ast.ActionBlock {
	action := &ast.ActionBlock{}
	p.nextToken()
	for !p.curTokenIs(token.RBRACE) {
		action.Actions = append(action.Actions, p.parseAction())
	}
	p.nextToken()
	return action
}

func (p *Parser) parseRegexAction() *ast.RegexAction {
	action := &ast.RegexAction{Rule: p.curToken.Literal}
	p.nextToken()
	action.Action = p.parseAction()
	return action
}

func (p *Parser) parseGotoAction() *ast.GotoAction {
	action := &ast.GotoAction{}
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		action.Target = p.curToken.Literal
	} else {
		p.GOTOsThatNeedNextStateAssigned = append(p.GOTOsThatNeedNextStateAssigned, action)
	}
	p.nextToken()
	return action
}

func (p *Parser) parseDoAction() *ast.DoSedAction {
	action := &ast.DoSedAction{Command: p.curToken.Literal}
	p.nextToken()
	return action
}

func (p *Parser) parsePrintAction() *ast.PrintAction {
	p.nextToken()
	return &ast.PrintAction{}
}
