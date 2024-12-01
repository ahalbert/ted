package flags

var Flags struct {
	ProgramFile string   `arg:"-f,--program-file" placeholder:"FSAFILE" help:"Finite State Autonoma file to run"`
	Program     string   `arg:"positional"`
	InputFiles  []string `arg:"positional" placeholder:"INPUTFILE"`
	NoPrint     bool     `arg:"-n,--no-print" help:"Do not print lines by default."`
	DebugMode   bool     `arg:"--debug" help:"Provides Lexer and Parser information"`
}
