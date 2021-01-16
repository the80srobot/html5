package html5

import (
	"fmt"
	"io"

	"github.com/the80srobot/html5/bindings"
)

// SubsectionNode represents a self-contained part of the page, which can be
// repeated in the output, each time with different bindings. (For example,
// comments under an article might each be a subsection.)
//
// Subsections can contain other subsections.
type SubsectionNode struct {
	Prototype Node
	Name      string
}

func (ns *SubsectionNode) Apply(n Node) error {
	switch n := n.(type) {
	case *ElementNode:
		n.Contents = append(n.Contents, ns)
	case *MultiNode:
		n.Contents = append(n.Contents, ns)
	default:
		return fmt.Errorf("SubsectionNode can only be applied to ElementNode or MultiNode, %v", n)
	}
	return nil
}

func (ns *SubsectionNode) compile(tc *templateCompiler, depth int, opts *CompileOptions) error {
	m := tc.bindings.Nest(ns.Name)
	subsectionOpts := *opts
	subsectionOpts.RootDepth = depth
	t, err := Compile(ns.Prototype, m, &subsectionOpts)
	if err != nil {
		return err
	}
	tc.appendChunk(subsectionChunk{template: *t, bindings: m})
	return nil
}

type subsectionChunk struct {
	template Template
	bindings *bindings.Map
}

func (sc subsectionChunk) build(w io.Writer, vm *bindings.ValueMap) error {
	stream := vm.GetStream(sc.bindings)
	if stream == nil {
		return nil
	}

	next := stream.Stream()
	for sectionValues := next(); sectionValues != nil; sectionValues = next() {
		if err := sc.template.GenerateHTML(w, sectionValues); err != nil {
			return err
		}
	}

	return nil
}

// func (sc subsectionChunk) String() string {
// 	var sb strings.Builder
// 	fmt.Fprintf(&sb, "subsection(%v) {\n", sc.bindings)

// 	scanner := bufio.NewScanner(strings.NewReader(sc.template.String()))
// 	scanner.Split(bufio.ScanLines)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		sb.WriteByte('\t')
// 		sb.WriteString(line)
// 		sb.WriteByte('\n')
// 	}

// 	sb.WriteString("\n}")
// 	return sb.String()
// }
