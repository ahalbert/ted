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
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
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
			tok = token.Token{Type: token.SLASH, Literal: "/"}
		} else {
			tok = token.Token{Type: token.REGEX, Literal: l.readUntilChar('/')}
		}
	case '"':
		l.readChar()
		tok = token.Token{Type: token.STRING, Literal: l.readUntilChar('"')}
	case '\'':
		l.readChar()
		tok = token.Token{Type: token.STRING, Literal: l.readUntilChar('\'')}
	case '`':
		l.readChar()
		tok = token.Token{Type: token.STRING, Literal: l.readUntilChar('`')}
	case '-':
		l.readChar()
		if l.ch == '-' && l.peek(1) == ">" {
			l.readChar()
			tok = token.Token{Type: token.RESET, Literal: "-->"}
		} else if l.ch == '>' {
			tok = token.Token{Type: token.GOTO, Literal: "->"}
		} else {
			tok = token.Token{Type: token.MINUS, Literal: "-"}
		}
	case '=':
		if l.peek(1) == "=" {
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: "=="}
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: "="}
		}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: "{"}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: "}"}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: "("}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: ")"}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: ","}
	case ':':
		tok = token.Token{Type: token.COLON, Literal: ":"}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: ";"}
	case '+':
		tok = token.Token{Type: token.PLUS, Literal: "+"}
	case '*':
		tok = token.Token{Type: token.ASTERISK, Literal: "*"}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
	default:
		if isLetter(l.ch) {
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
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()

	return tok
}

func (l *Lexer) handleIdentfierSpecialCases(t token.Token) token.Token {
	if l.ch == ':' {
		l.readChar()
		return token.Token{Type: token.LABEL, Literal: t.Literal}
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
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
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

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
