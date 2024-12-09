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

	//Keywords
	DO      = "DO"
	START   = "START"
	STOP    = "STOP"
	CAPTURE = "CAPTURE"
	LABEL   = "LABEL"
	PRINT   = "PRINT"
	CLEAR   = "CLEAR"
)

var keywords = map[string]TokenType{
	"do":      DO,
	"capture": CAPTURE,
	"print":   PRINT,
	"start":   START,
	"stop":    STOP,
	"clear":   CLEAR,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
