package lexer

import (
	"slices"

	"github.com/ahalbert/ted/ted/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	lineNum      int
	linePosition int
}

func New(input string) *Lexer {
	l := &Lexer{input: input, lineNum: 1, linePosition: 1}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]

		if l.ch == '\n' {
			l.lineNum++
			l.linePosition = 1
		}
	}
	l.linePosition++
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peek(lookahead int) string {
	if l.readPosition+lookahead >= len(l.input) {
		return ""
	} else {
		return l.input[l.readPosition : l.readPosition+lookahead]
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '/':
		l.readChar()
		if l.ch == ' ' {
			tok = l.newToken(token.SLASH, "/")
		} else {
			tok = l.newToken(token.REGEX, l.readUntilChar('/'))
		}
	case '"':
		l.readChar()
		tok = l.newToken(token.STRING, l.readUntilChar('"'))
	case '\'':
		l.readChar()
		tok = l.newToken(token.STRING, l.readUntilChar('\''))
	case '`':
		l.readChar()
		tok = l.newToken(token.STRING, l.readUntilChar('`'))
	case '-':
		l.readChar()
		if l.ch == '-' && l.peek(1) == ">" {
			l.readChar()
			tok = l.newToken(token.RESET, "-->")
		} else if l.ch == '>' {
			tok = l.newToken(token.GOTO, "->")
		} else {
			tok = l.newToken(token.MINUS, "-")
		}
	case '=':
		if l.peek(1) == "=" {
			l.readChar()
			tok = l.newToken(token.EQ, "==")
		} else {
			tok = l.newToken(token.ASSIGN, "=")
		}
	case '{':
		tok = l.newToken(token.LBRACE, "{")
	case '}':
		tok = l.newToken(token.RBRACE, "}")
	case '(':
		tok = l.newToken(token.LPAREN, "(")
	case ')':
		tok = l.newToken(token.RPAREN, ")")
	case ',':
		tok = l.newToken(token.COMMA, ",")
	case ':':
		tok = l.newToken(token.COLON, ":")
	case ';':
		tok = l.newToken(token.SEMICOLON, ";")
	case '+':
		tok = l.newToken(token.PLUS, "+")
	case '*':
		tok = l.newToken(token.ASTERISK, "*")
	case 0:
		tok = l.newToken(token.EOF, "")
	default:
		if isLetter(l.ch) {
			tok.LineNum = l.lineNum
			tok.Position = l.linePosition
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			switch tok.Type {
			case token.DO:
				tok.Literal = l.readDo()
				l.readChar()
			case token.DOUNTIL:
				tok.Literal = l.readDo()
				l.readChar()
			case token.IDENT:
				tok = l.handleIdentfierSpecialCases(tok)
			}
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL, string(l.ch))
		}
	}

	l.readChar()

	return tok
}

func (l *Lexer) handleIdentfierSpecialCases(t token.Token) token.Token {
	if l.ch == ':' {
		l.readChar()
		return l.newToken(token.LABEL, t.Literal)
	}
	return t
}

func (l *Lexer) readDo() string {
	l.skipWhitespace()
	switch l.ch {
	case '"':
		l.readChar()
		return l.readUntilChar('"')
	case '\'':
		l.readChar()
		return l.readUntilChar('\'')
	case '`':
		l.readChar()
		return l.readUntilChar('`')
	default:
		return l.readUntilChar(' ', '\t', '\n', '\r')
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) && l.ch != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == '#' || l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '#' {
			l.readUntilChar('\n')
		}
		l.readChar()
	}
}

func (l *Lexer) readUntilChar(chars ...byte) string {
	position := l.position
	for !slices.Contains(chars, l.ch) && l.ch != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || '0' <= ch && ch <= '9' || ch == '_' || ch == '$'
}

func (l *Lexer) newToken(tokenType token.TokenType, s string) token.Token {
	return token.Token{Type: tokenType, Literal: s, LineNum: l.lineNum, Position: l.linePosition}
}
