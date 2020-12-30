package html

import (
	"fmt"
	"io"
)

type Case struct {
	Condition Condition
	Output    Node
}

type Condition func(*ValueSet) bool

type SwitchNode struct {
	Cases   []Case
	Default Node
}

func (sn *SwitchNode) compile(tc *templateCompiler, depth int, opts *CompileOptions) error {
	sc := switchChunk{
		conditions: make([]Condition, len(sn.Cases)),
		templates:  make([]*Template, len(sn.Cases)+1),
	}

	for i, c := range sn.Cases {
		sc.conditions[i] = c.Condition
		if c.Output == nil {
			continue
		}
		t, err := Compile(c.Output, depth, opts)
		if err != nil {
			return fmt.Errorf("compiling case %d/%d: %w", i+1, len(sn.Cases), err)
		}
		sc.templates[i] = t
	}

	if sn.Default != nil {
		t, err := Compile(sn.Default, depth, opts)
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

func (sc switchChunk) build(w io.Writer, vs *ValueSet) error {
	for i, c := range sc.conditions {
		if c(vs) {
			if t := sc.templates[i]; t != nil {
				return t.GenerateHTML(w, vs)
			}
			return nil
		}
	}

	if t := sc.templates[len(sc.templates)-1]; t != nil {
		return t.GenerateHTML(w, vs)
	}
	return nil
}

// func (sc componentChunk) String() string {
// 	var sb strings.Builder
// 	fmt.Fprint(&sb, "switch {\n")

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
