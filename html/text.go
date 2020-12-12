package html

import (
	"bytes"
	"fmt"
	"io"
)

type TextNode struct {
	Constant    SafeString
	StringName  string
	Width       int
	IndentStyle IndentStyle
}

func (t *TextNode) String() string {
	return fmt.Sprintf("&TextNode{constant=%q, binding=%v, width=%d, style=%v}", t.Constant, t.StringName, t.Width, t.IndentStyle)
}

func (t *TextNode) compile(tc *templateCompiler, depth int, opts *CompileOptions) error {
	is := t.IndentStyle
	if opts.Compact {
		is = Inline
	}
	return appendText(tc, depth, t, is, opts.Indent)
}

type textBindingChunk struct {
	TextNode
	stringTag   Tag
	indent      string
	depth       int
	indentStyle IndentStyle
}

func (tc textBindingChunk) build(w io.Writer, vs *ValueSet) error {
	// An optimization: if we don't need to break up lines or indent anything,
	// then we can just print the binding value as is.
	if tc.indentStyle == Inline && tc.Width <= 0 {
		return vs.writeStringTo(w, tc.stringTag)
	}

	var b bytes.Buffer
	if err := vs.writeStringTo(&b, tc.stringTag); err != nil {
		return err
	}
	return fprintBlockText(w, tc.depth, tc.Width, tc.indent, tc.indentStyle, &b)
}

func (tc textBindingChunk) String() string {
	return fmt.Sprintf("textBinding{%v, tag=%v, indent=%q, depth=%d, style=%v}",
		&tc.TextNode, tc.stringTag, tc.indent, tc.depth, tc.indentStyle)
}
