package html

import (
	"fmt"
	"io"
	"strings"
)

type Node interface {
	compile(tc *templateCompiler, depth int, opts *CompileOptions) error
}

type Template struct {
	chunks   []chunk
	Bingings *BindingSet
}

func (t *Template) GenerateHTML(w io.Writer, vs *ValueSet) error {
	for i, chunk := range t.chunks {
		if err := chunk.build(w, vs); err != nil {
			return fmt.Errorf("building chunk #%d of %d: %w", i, len(t.chunks), err)
		}
	}
	return nil
}

func (t *Template) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Template{\n")
	for name, tag := range t.Bingings.stringNames {
		fmt.Fprintf(&sb, "  binding tag %s: %d\n", name, tag)
	}
	for i, c := range t.chunks {
		fmt.Fprintf(&sb, "  chunk %d/%d: %v\n", i+1, len(t.chunks), c)
	}
	sb.WriteByte('}')
	return sb.String()
}

func Compile(n Node, depth int, bs *BindingSet, opts *CompileOptions) (*Template, error) {
	tc := &templateCompiler{bindings: bs}
	tc.separateChunks = opts.SeparateStaticChunks
	if err := n.compile(tc, depth, opts); err != nil {
		return nil, err
	}
	tc.flush()
	return &Template{chunks: tc.chunks, Bingings: tc.bindings}, nil
}

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
}

func (opts *CompileOptions) String() string {
	return fmt.Sprintf("{Indent: %q, Compact: %v, SeparateStaticChunks: %v}", opts.Indent, opts.Compact, opts.SeparateStaticChunks)
}

type chunk interface {
	build(w io.Writer, vs *ValueSet) error
}

type staticChunk struct {
	data string
}

func (sc staticChunk) build(w io.Writer, _ *ValueSet) error {
	_, err := io.WriteString(w, sc.data)
	return err
}

func (sc staticChunk) String() string {
	return fmt.Sprintf("static{%q}", sc.data)
}

type stringBindingChunk struct {
	stringTag Tag
}

func (sbc stringBindingChunk) build(w io.Writer, vs *ValueSet) error {
	return vs.writeStringTo(w, sbc.stringTag)
}

func (sbc stringBindingChunk) String() string {
	return fmt.Sprintf("stringBinding{%v}", sbc.stringTag)
}
