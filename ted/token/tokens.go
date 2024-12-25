package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	//Identfiers
	IDENT  = "IDENT"
	REGEX  = "REGEX"
	STRING = "STRING"

	//symbols

	COLON     = ":"
	ASSIGN    = "="
	SEMICOLON = ";"
	COMMA     = ","
	GOTO      = "->"
	RESET     = "-->"
	LBRACE    = "{"
	RBRACE    = "}"

	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	LPAREN = "("
	RPAREN = ")"

	//Keywords
	DO       = "DO"
	DOUNTIL  = "DOUNTIL"
	START    = "START"
	STOP     = "STOP"
	CAPTURE  = "CAPTURE"
	LABEL    = "LABEL"
	LET      = "LET"
	PRINT    = "PRINT"
	PRINTLN  = "PRINTLN"
	CLEAR    = "CLEAR"
	REWIND   = "REWIND"
	FASTFWD  = "FASTFORWARD"
	PAUSE    = "PAUSE"
	PLAY     = "PLAY"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	FUNCTION = "FUNCTION"
)

var keywords = map[string]TokenType{
	"do":          DO,
	"dountil":     DOUNTIL,
	"capture":     CAPTURE,
	"print":       PRINT,
	"println":     PRINTLN,
	"start":       START,
	"stop":        STOP,
	"clear":       CLEAR,
	"let":         LET,
	"rewind":      REWIND,
	"fastforward": FASTFWD,
	"pause":       PAUSE,
	"play":        PLAY,
	"if":          IF,
	"else":        ELSE,
	"function":    FUNCTION,
	"return":      RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
