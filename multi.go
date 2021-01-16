package html5

import "fmt"

// MultiNode concatenates several other nodes.
type MultiNode struct {
	Contents []Node
}

func Multi(contents ...Content) *MultiNode {
	m := &MultiNode{}
	for _, c := range contents {
		c.Apply(m)
	}
	return m
}

// Apply will insert this node as a child into the other node.
func (m *MultiNode) Apply(n Node) error {
	switch n := n.(type) {
	case *ElementNode:
		n.Contents = append(n.Contents, m)
	case *MultiNode:
		n.Contents = append(n.Contents, m.Contents...)
	default:
		return fmt.Errorf("MultiNode can only be applied to ElementNode or MultiNode, got %v", n)
	}

	return nil
}

func (m *MultiNode) compile(db *templateCompiler, depth int, opts *CompileOptions) error {
	for _, c := range m.Contents {
		if err := c.compile(db, depth, opts); err != nil {
			return err
		}
	}
	return nil
}
