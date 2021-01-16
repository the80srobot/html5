package html

import "fmt"

type Content interface {
	Apply(n Node) error
}

type NodeOption func(n Node) error

func (f NodeOption) Apply(n Node) error {
	return f(n)
}

func Indent(is IndentStyle) Content {
	f := func(n Node) error {
		switch n := n.(type) {
		case *ElementNode:
			n.IndentStyle = is
		default:
			return fmt.Errorf("Indent option cannot be applied to node %v", n)
		}
		return nil
	}
	return NodeOption(f)
}

func SelfClosing() Content {
	f := func(n Node) error {
		switch n := n.(type) {
		case *ElementNode:
			n.SelfClosing = true
		default:
			return fmt.Errorf("SelfClosing option cannot be applied to node %v", n)
		}
		return nil
	}
	return NodeOption(f)
}
