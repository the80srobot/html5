package html

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/the80srobot/html5/bindings"
	"github.com/the80srobot/html5/safe"
)

type TextNode struct {
	Value Value
	Width int
}

func (t *TextNode) String() string {
	return fmt.Sprintf("&TextNode{value=%v, width=%d}", t.Value, t.Width)
}

func (t *TextNode) compile(tc *templateCompiler, depth int, opts *CompileOptions) error {
	switch v := t.Value.(type) {
	case safe.String:
		s, err := safe.Check(v, safe.TextSafe)
		if err != nil {
			return err
		}
		return fprintBlockText(tc, depth, t.Width, opts.Indent, strings.NewReader(s))
	case bindings.Var:
		v = tc.bindings.Attach(v, safe.TextSafe)
		tc.appendChunk(textBindingChunk{TextNode: *t, depth: depth, indent: opts.Indent, binding: v})
		return nil
	default:
		return fmt.Errorf("value must be safe.String or *bindings.Var, %v (%v) is neither", v, reflect.TypeOf(v))
	}
}

type textBindingChunk struct {
	TextNode
	binding bindings.Var
	indent  string
	depth   int
}

func (tc textBindingChunk) build(w io.Writer, vm *bindings.ValueMap) error {
	// An optimization: if we don't need to break up the lines then we can just
	// print the binding value as is.
	if tc.Width <= 0 {
		_, err := io.WriteString(w, vm.GetString(tc.binding))
		return err
	}

	var b bytes.Buffer
	if _, err := io.WriteString(&b, vm.GetString(tc.binding)); err != nil {
		return err
	}
	return fprintBlockText(w, tc.depth, tc.Width, tc.indent, &b)
}

func (tc textBindingChunk) String() string {
	return fmt.Sprintf("textBinding{%v, tag=%v, indent=%q, depth=%d}",
		&tc.TextNode, tc.binding, tc.indent, tc.depth)
}
