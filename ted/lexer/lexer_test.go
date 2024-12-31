package lexer

import (
	"testing"

	"github.com/ahalbert/ted/ted/token"
)

func TestReadChar(t *testing.T) {
	tests := []struct {
		input                 string
		expectedChars         []byte
		expectedPositions     []int
		expectedReadPositions []int
	}{
		{
			input:                 "test",
			expectedChars:         []byte{'t', 'e', 's', 't', 0},
			expectedPositions:     []int{0, 1, 2, 3, 4},
			expectedReadPositions: []int{1, 2, 3, 4, 5},
		},
		{
			input:                 "hello",
			expectedChars:         []byte{'h', 'e', 'l', 'l', 'o', 0},
			expectedPositions:     []int{0, 1, 2, 3, 4, 5},
			expectedReadPositions: []int{1, 2, 3, 4, 5, 6},
		},
		{
			input:                 "a\nb",
			expectedChars:         []byte{'a', '\n', 'b', 0},
			expectedPositions:     []int{0, 1, 2, 3},
			expectedReadPositions: []int{1, 2, 3, 4},
		},
		{
			input:                 "",
			expectedChars:         []byte{0},
			expectedPositions:     []int{0},
			expectedReadPositions: []int{1},
		},
	}

	for i, tt := range tests {
		l := New(tt.input)
		for j := 0; j < len(tt.expectedChars); j++ {
			if l.ch != tt.expectedChars[j] {
				t.Errorf("test[%d] - readChar() ch = %q, want %q", i, l.ch, tt.expectedChars[j])
			}
			if l.position != tt.expectedPositions[j] {
				t.Errorf("test[%d] - readChar() position = %d, want %d", i, l.position, tt.expectedPositions[j])
			}
			if l.readPosition != tt.expectedReadPositions[j] {
				t.Errorf("test[%d] - readChar() readPosition = %d, want %d", i, l.readPosition, tt.expectedReadPositions[j])
			}
			l.readChar()
		}
	}
}

func TestPeek(t *testing.T) {
	input := "test"
	l := New(input)

	tests := []struct {
		lookahead    int
		expectedPeek string
	}{
		{0, ""},
		{1, "e"},
		{2, "es"},
		{3, ""},
		{4, ""},
		{5, ""},
	}

	for i, tt := range tests {
		peeked := l.peek(tt.lookahead)
		if peeked != tt.expectedPeek {
			t.Errorf("test[%d] - peek(%d) = %q, want %q", i, tt.lookahead, peeked, tt.expectedPeek)
		}
	}
}

func TestReadDo(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{`"hello w"`, "hello w"},
		{`' h world'`, " h world"},
		{"`test`", "test"},
		{"simple text", "simple"},
		{"  whitespace", "whitespace"},
		{"whitespace  ", "whitespace"},
		{"#whitespace \nab", "ab"},
	}

	for i, tt := range tests {
		l := New(tt.input)
		output := l.readDo()
		if output != tt.expectedOutput {
			t.Errorf("test[%d] - readDo() = %q, want %q", i, output, tt.expectedOutput)
		}
	}
}

func TestReadIdentifier(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"identifier", "identifier"},
		{"_underscore", "_underscore"},
		{"$dollar", "$dollar"},
		{"with123numbers", "with123numbers"},
		{"059mixed_Case$", "059mixed_Case$"},
		{"", ""},
		{"!notIdentifier", ""},
	}

	for i, tt := range tests {
		l := New(tt.input)
		output := l.readIdentifier()
		if output != tt.expectedOutput {
			t.Errorf("test[%d] - readIdentifier() = %q, want %q", i, output, tt.expectedOutput)
		}
	}
}

func TestNextToken(t *testing.T) {
	tests := []struct {
		input          string
		expectedTokens []struct {
			expectedType    token.TokenType
			expectedLiteral string
		}
	}{
		{
			input: `=+(){},;*:`, // char tests
			expectedTokens: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.ASSIGN, "="},
				{token.PLUS, "+"},
				{token.LPAREN, "("},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.RBRACE, "}"},
				{token.COMMA, ","},
				{token.SEMICOLON, ";"},
				{token.ASTERISK, "*"},
				{token.COLON, ":"},
				{token.EOF, ""},
			},
		},
		{
			input: `/ab{2}/`, // regexp test
			expectedTokens: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.REGEX, "ab{2}"},
				{token.EOF, ""},
			},
		},
		{
			input: `/ +`, // slash test
			expectedTokens: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.SLASH, "/"},
				{token.PLUS, "+"},
				{token.EOF, ""},
			},
		},
		{
			input: `"abcd"`, // string test
			expectedTokens: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.STRING, "abcd"},
				{token.EOF, ""},
			},
		},
		{
			input: `'abcd'`, // string test
			expectedTokens: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.STRING, "abcd"},
				{token.EOF, ""},
			},
		},
		{
			input: "`abcd`", // string test
			expectedTokens: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.STRING, "abcd"},
				{token.EOF, ""},
			},
		},
		{
			input: `!`, // illegal char test
			expectedTokens: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.ILLEGAL, "!"},
				{token.EOF, ""},
			},
		},
		{
			input: `/foo/ -> /bar/ -> do s/baz/bang/`, // sample input
			expectedTokens: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.REGEX, "foo"},
				{token.GOTO, "->"},
				{token.REGEX, "bar"},
				{token.GOTO, "->"},
				{token.DO, "s/baz/bang/"},
				{token.EOF, ""},
			},
		},
		{
			input: `{ capture fastforward /buzz/ -> }`, // sample input
			expectedTokens: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.LBRACE, "{"},
				{token.CAPTURE, "capture"},
				{token.FASTFWD, "fastforward"},
				{token.REGEX, "buzz"},
				{token.GOTO, "->"},
				{token.RBRACE, "}"},
				{token.EOF, ""},
			},
		},
		{
			input: `dountil s/buzz/boop/ ->`, // sample input
			expectedTokens: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.DOUNTIL, "s/buzz/boop/"},
				{token.GOTO, "->"},
				{token.EOF, ""},
			},
		},
		{
			input: `# Welcome to Ted!
			/foo/ -> /bar/ -> do s/baz/bang/`, // sample input
			expectedTokens: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.REGEX, "foo"},
				{token.GOTO, "->"},
				{token.REGEX, "bar"},
				{token.GOTO, "->"},
				{token.DO, "s/baz/bang/"},
				{token.EOF, ""},
			},
		},
		{
			input: `/buzz/ {println myvar}`, // sample input
			expectedTokens: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.REGEX, "buzz"},
				{token.LBRACE, "{"},
				{token.PRINTLN, "println"},
				{token.IDENT, "myvar"},
				{token.RBRACE, "}"},
				{token.EOF, ""},
			},
		},
		{
			input: `ALL: /Success/ --> `, // sample input
			expectedTokens: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.LABEL, "ALL"},
				{token.REGEX, "Success"},
				{token.RESET, "-->"},
				{token.EOF, ""},
			},
		},
		{
			input: `5: { println count /div/ let count = count - 1 if count == 0 -> }`, // sample input
			expectedTokens: []struct {
				expectedType    token.TokenType
				expectedLiteral string
			}{
				{token.LABEL, "5"},
				{token.LBRACE, "{"},
				{token.PRINTLN, "println"},
				{token.IDENT, "count"},
				{token.REGEX, "div"},
				{token.LET, "let"},
				{token.IDENT, "count"},
				{token.ASSIGN, "="},
				{token.IDENT, "count"},
				{token.MINUS, "-"},
				{token.IDENT, "1"},
				{token.IF, "if"},
				{token.IDENT, "count"},
				{token.EQ, "=="},
				{token.IDENT, "0"},
				{token.GOTO, "->"},
				{token.RBRACE, "}"},
				{token.EOF, ""},
			},
		},
	}

	for _, tt := range tests {
		l := New(tt.input)

		for i, expectedToken := range tt.expectedTokens {
			tok := l.NextToken()

			if tok.Type != expectedToken.expectedType {
				t.Fatalf("input: %q, tests[%d] - tokentype wrong. expected=%q, got=%q", tt.input, i, expectedToken.expectedType, tok.Type)
			}

			if tok.Literal != expectedToken.expectedLiteral {
				t.Fatalf("input: %q, tests[%d] - literal wrong. expected=%q, got=%q", tt.input, i, expectedToken.expectedLiteral, tok.Literal)
			}
		}
	}
}
