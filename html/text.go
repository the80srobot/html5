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
}

func Text(contents ...Value) Node {
	switch len(contents) {
	case 0:
		return &TextNode{}
	case 1:
		return &TextNode{Value: contents[0]}
	default:
		m := &MultiNode{Contents: make([]Node, 0, len(contents))}
		for _, c := range contents {
			m.Contents = append(m.Contents, &TextNode{Value: c})
		}
		return m
	}
}

func (t *TextNode) Apply(n Node) error {
	switch n := n.(type) {
	case *ElementNode:
		n.Contents = append(n.Contents, t)
	case *MultiNode:
		n.Contents = append(n.Contents, t)
	default:
		return fmt.Errorf("TextNode can only be applied to ElementNode or MultiNode, got %v", n)
	}
	return nil
}

func (t *TextNode) String() string {
	return fmt.Sprintf("&TextNode{value=%v}", t.Value)
}

func (t *TextNode) compile(tc *templateCompiler, depth int, opts *CompileOptions) error {
	switch v := t.Value.(type) {
	case safe.String:
		s, err := safe.Check(v, safe.TextSafe)
		if err != nil {
			return err
		}
		return fprintBlockText(tc, depth, opts.TextWidth, opts.Indent, strings.NewReader(s))
	case bindings.Var:
		v = tc.bindings.Attach(v, safe.TextSafe)
		tc.appendChunk(textBindingChunk{TextNode: *t, depth: depth, indent: opts.Indent, binding: v, width: opts.TextWidth})
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
	width   int
}

func (tc textBindingChunk) build(w io.Writer, vm *bindings.ValueMap) error {
	// An optimization: if we don't need to break up the lines then we can just
	// print the binding value as is.
	if tc.width <= 0 {
		_, err := io.WriteString(w, vm.GetString(tc.binding))
		return err
	}

	var b bytes.Buffer
	if _, err := io.WriteString(&b, vm.GetString(tc.binding)); err != nil {
		return err
	}
	return fprintBlockText(w, tc.depth, tc.width, tc.indent, &b)
}

func (tc textBindingChunk) String() string {
	return fmt.Sprintf("textBinding{%v, tag=%v, indent=%q, depth=%d}",
		&tc.TextNode, tc.binding, tc.indent, tc.depth)
}
