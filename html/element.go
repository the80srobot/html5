package html

import (
	"errors"
	"fmt"
)

// ElementNode represents an HTML element, like <p>.
type ElementNode struct {
	Name                string
	Attributes          []Attribute
	Contents            []Node
	IndentStyle         IndentStyle
	SelfClosing         bool
	XMLStyleSelfClosing bool
}

func Element(name string, contents ...Content) *ElementNode {
	e := &ElementNode{Name: name}
	for _, c := range contents {
		c.Apply(e)
	}
	return e
}

func (e *ElementNode) Apply(n Node) error {
	switch n := n.(type) {
	case *ElementNode:
		n.Contents = append(n.Contents, e)
	case *MultiNode:
		n.Contents = append(n.Contents, e)
	default:
		return fmt.Errorf("ElementNode can only be applied to ElementNode or MultiNode, got %v", n)
	}
	return nil
}

func (e *ElementNode) deduplicateAttributes() {
	var attrs []Attribute
	seen := make(map[string]struct{}, len(e.Attributes))
	for i := len(e.Attributes) - 1; i >= 0; i-- {
		if _, ok := seen[e.Attributes[i].Name]; ok {
			continue
		}
		seen[e.Attributes[i].Name] = struct{}{}
		attrs = append(attrs, e.Attributes[i])
	}

	// attrs are deduplicated and also reversed.
	l := len(attrs)
	for i := 0; i < l/2; i++ {
		attrs[i], attrs[l-i-1] = attrs[l-i-1], attrs[i]
	}
	e.Attributes = attrs
}

func (e *ElementNode) compile(tc *templateCompiler, depth int, opts *CompileOptions) error {
	e.deduplicateAttributes()

	isBlock := e.IndentStyle == Block && !opts.Compact

	// Block elements always start on a new line.
	if isBlock && !tc.freshLine() {
		if err := tc.appendNewLine(depth, opts.Indent); err != nil {
			return err
		}
	}

	// Handle the xHTML style (produce <br /> instead of <br>).
	openingTag := tagOpen
	if e.XMLStyleSelfClosing {
		openingTag = tagSelfClose
	}

	if err := appendTag(tc, e.Name, openingTag, e.Attributes...); err != nil {
		return err
	}

	if e.SelfClosing || e.XMLStyleSelfClosing {
		if len(e.Contents) != 0 {
			return errors.New("self-closing element cannot have contents")
		}
		return nil
	}

	if isBlock {
		depth++
		// Block element contents are indented by an additional level.
		if err := tc.appendNewLine(depth, opts.Indent); err != nil {
			return err
		}
	}

	for _, c := range e.Contents {
		if err := c.compile(tc, depth, opts); err != nil {
			return err
		}
	}

	if isBlock {
		depth--
		if err := tc.appendNewLine(depth, opts.Indent); err != nil {
			return err
		}
	}

	if err := appendTag(tc, e.Name, tagClose); err != nil {
		return err
	}

	return nil
}
