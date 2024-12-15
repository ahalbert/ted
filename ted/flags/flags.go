package flags

var Flags struct {
	ProgramFile string   `arg:"-f,--fsa-file" placeholder:"FSAFILE" help:"Finite State Autonoma file to run."`
	NoPrint     bool     `arg:"-n,--no-print" help:"Do not print lines by default."`
	DebugMode   bool     `arg:"--debug" help:"Provides Lexer and Parser information."`
	Variables   []string `arg:"--var,separate" placeholder:"key=value" help:"Variable in the format name=value."`
	Program     string   `arg:"positional" help:"Program to run."`
	InputFiles  []string `arg:"positional" placeholder:"INPUTFILE" help:"File to use as input."`
}
