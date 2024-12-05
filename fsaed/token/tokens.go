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
	SEMICOLON = ";"
	COMMA     = ","
	GOTO      = "->"
	RESET     = "-->"
	LBRACE    = "{"
	RBRACE    = "}"

	//Keywords
	DO    = "DO"
	LABEL = "LABEL"
	PRINT = "PRINT"
	START = "START"
	STOP  = "STOP"
)

var keywords = map[string]TokenType{
	"do":    DO,
	"print": PRINT,
	"start": START,
	"stop":  STOP,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
