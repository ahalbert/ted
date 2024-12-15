package lexer

import (
	"testing"
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
