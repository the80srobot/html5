package html

import (
	"bufio"
	"fmt"
	"io"
	"strings"
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

func (ns *SubsectionNode) compile(tc *templateCompiler, depth int, opts *CompileOptions) error {
	t, err := Compile(ns.Prototype, depth, opts)
	if err != nil {
		return err
	}
	tag := tc.bindings.DeclareSubsection(ns.Name, &t.Bingings)
	tc.appendChunk(subsectionChunk{template: *t, binding: tag})
	return nil
}

type subsectionChunk struct {
	template Template
	binding  Tag
}

func (sc subsectionChunk) build(w io.Writer, vs *ValueSet) error {
	next, err := vs.iterateSubsection(sc.binding)
	if err != nil {
		return err
	}

	for sectionValues := next(); sectionValues != nil; sectionValues = next() {
		if err := sc.template.GenerateHTML(w, sectionValues); err != nil {
			return err
		}
	}

	return nil
}

func (sc subsectionChunk) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "subsection(%v) {\n", sc.binding)

	scanner := bufio.NewScanner(strings.NewReader(sc.template.String()))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		sb.WriteByte('\t')
		sb.WriteString(line)
		sb.WriteByte('\n')
	}

	sb.WriteString("\n}")
	return sb.String()
}
