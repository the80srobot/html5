package html5

import (
	"fmt"

	"github.com/the80srobot/html/html5/html"
)

type Option interface {
	Apply(n html.Node) error
}

type funcOption func(n html.Node) error

func (f funcOption) Apply(n html.Node) error {
	return f(n)
}

func Indent(is html.IndentStyle) Option {
	f := func(n html.Node) error {
		switch n := n.(type) {
		case *html.ElementNode:
			n.IndentStyle = is
		case *html.TextNode:
			n.IndentStyle = is
		default:
			return fmt.Errorf("Indent option cannot be applied to node %v", n)
		}
		return nil
	}
	return funcOption(f)
}

func LineWidth(width int) Option {
	f := func(n html.Node) error {
		t, ok := n.(*html.TextNode)
		if !ok {
			return fmt.Errorf("LineWidth option can only be applied to text nodes, got %v", n)
		}
		t.Width = width
		return nil
	}
	return funcOption(f)
}

func applyOptions(n html.Node, opts ...Option) {
	for _, opt := range opts {
		if err := opt.Apply(n); err != nil {
			panic(fmt.Sprintf("Invalid option (programmer error): %v", err))
		}
	}
}
