package html

import (
	"bytes"
	"fmt"
	"io"
)

type TextNode struct {
	Value SafeString
	Width int
}

func (t *TextNode) String() string {
	return fmt.Sprintf("&TextNode{value=%v, width=%d}", t.Value, t.Width)
}

func (t *TextNode) compile(tc *templateCompiler, depth int, opts *CompileOptions) error {
	return appendText(tc, depth, t, opts.Indent)
}

type textBindingChunk struct {
	TextNode
	stringTag Tag
	indent    string
	depth     int
}

func (tc textBindingChunk) build(w io.Writer, vs *ValueSet) error {
	// An optimization: if we don't need to break up the lines then we can just
	// print the binding value as is.
	if tc.Width <= 0 {
		return vs.writeStringTo(w, tc.stringTag)
	}

	var b bytes.Buffer
	if err := vs.writeStringTo(&b, tc.stringTag); err != nil {
		return err
	}
	return fprintBlockText(w, tc.depth, tc.Width, tc.indent, &b)
}

func (tc textBindingChunk) String() string {
	return fmt.Sprintf("textBinding{%v, tag=%v, indent=%q, depth=%d}",
		&tc.TextNode, tc.stringTag, tc.indent, tc.depth)
}
