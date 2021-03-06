package html5

import (
	"fmt"
	"io"

	"github.com/the80srobot/html5/bindings"
)

type Case struct {
	Condition Condition
	Output    Node
}

type Condition func(*bindings.ValueMap) bool

type SwitchNode struct {
	Cases   []Case
	Default Node
}

func (sn *SwitchNode) Apply(n Node) error {
	switch n := n.(type) {
	case *ElementNode:
		n.Contents = append(n.Contents, sn)
	case *MultiNode:
		n.Contents = append(n.Contents, sn)
	default:
		return fmt.Errorf("SwitchNode can only be applied to ElementNode or MultiNode, %v", n)
	}
	return nil
}

func (sn *SwitchNode) compile(tc *templateCompiler, depth int, opts *CompileOptions) error {
	sc := switchChunk{
		conditions: make([]Condition, len(sn.Cases)),
		templates:  make([]*Template, len(sn.Cases)+1),
	}

	nestedOpts := *opts
	nestedOpts.RootDepth = depth
	for i, c := range sn.Cases {
		sc.conditions[i] = c.Condition
		if c.Output == nil {
			continue
		}
		t, err := Compile(c.Output, tc.bindings, &nestedOpts)
		if err != nil {
			return fmt.Errorf("compiling case %d/%d: %w", i+1, len(sn.Cases), err)
		}
		sc.templates[i] = t
	}

	if sn.Default != nil {
		t, err := Compile(sn.Default, tc.bindings, &nestedOpts)
		if err != nil {
			return fmt.Errorf("compiling default case: %w", err)
		}
		sc.templates[len(sc.templates)-1] = t
	}

	tc.appendChunk(sc)
	return nil
}

type switchChunk struct {
	conditions []Condition
	templates  []*Template
}

func (sc switchChunk) build(w io.Writer, vm *bindings.ValueMap) error {
	for i, c := range sc.conditions {
		if c(vm) {
			if t := sc.templates[i]; t != nil {
				return t.GenerateHTML(w, vm)
			}
			return nil
		}
	}

	if t := sc.templates[len(sc.templates)-1]; t != nil {
		return t.GenerateHTML(w, vm)
	}
	return nil
}
