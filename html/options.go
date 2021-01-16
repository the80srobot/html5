package html

import "fmt"

type IndentStyle int16

const (
	Inline IndentStyle = iota
	Block
)

var (
	Compact = CompileOptions{
		Indent:  "  ",
		Compact: true,
	}

	Tidy = CompileOptions{
		Indent:  "  ",
		Compact: false,
	}

	Debug = CompileOptions{
		Indent:               "  ",
		Compact:              false,
		SeparateStaticChunks: true,
	}
)

type CompileOptions struct {
	Indent               string
	Compact              bool
	SeparateStaticChunks bool
	TextWidth            int
}

func (opts *CompileOptions) String() string {
	return fmt.Sprintf("{Indent: %q, Compact: %v, SeparateStaticChunks: %v}", opts.Indent, opts.Compact, opts.SeparateStaticChunks)
}
