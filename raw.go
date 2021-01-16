package html5

import (
	"fmt"
	"reflect"

	"github.com/the80srobot/html5/bindings"
	"github.com/the80srobot/html5/safe"
)

// RawNode inserts a fully trusted string directly into the page.
type RawNode struct {
	HTML Value
}

func (r *RawNode) Apply(n Node) error {
	switch n := n.(type) {
	case *ElementNode:
		n.Contents = append(n.Contents, r)
	case *MultiNode:
		n.Contents = append(n.Contents, r)
	default:
		return fmt.Errorf("RawNode can only be applied to ElementNode or MultiNode, got %v", n)
	}
	return nil
}

func (r *RawNode) compile(tc *templateCompiler, _ int, _ *CompileOptions) error {
	switch v := r.HTML.(type) {
	case safe.String:
		s, err := safe.Check(v, safe.HTMLSafe)
		if err != nil {
			return err
		}
		_, err = tc.WriteString(s)
		return err
	case bindings.Var:
		tc.appendVar(v, safe.HTMLSafe)
		return nil
	default:
		return fmt.Errorf("value must be safe.String or *bindings.Var, %v (%v) is neither", v, reflect.TypeOf(v))
	}
}
