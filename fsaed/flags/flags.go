package flags

var Flags struct {
	Program string `arg:"positional"`
	NoPrint bool   `arg:"-n,--no-print" help:"Do not print lines by default."`
}
