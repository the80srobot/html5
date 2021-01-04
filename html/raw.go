package html

import (
	"fmt"

	"github.com/the80srobot/html5/bindings"
	"github.com/the80srobot/html5/safe"
)

// RawNode inserts a fully trusted string directly into the page.
type RawNode struct {
	HTML Value
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
		return fmt.Errorf("value must be safe.String or *bindings.Var, %v is neither", v)
	}
}
