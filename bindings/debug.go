package bindings

import (
	"io"
	"strings"
)

// DebugDumper knows how to print a verbose description of itself to a provided
// writer. The depth argument controls indentation.
type DebugDumper interface {
	DebugDump(w io.Writer, depth int)
}

// DebugString converts a DebugDumper to a verbose string representation for
// debugging or human-readable diffing.
//
// In particular, this is useful with cmp.Transformer to generate human-readable
// diffs in tests. For example:
//
//  cmp.Transformer(func(vm *ValueMap) string { return DebugString(vm) })
func DebugString(dd DebugDumper) string {
	var sb strings.Builder
	dd.DebugDump(&sb, 0)
	return sb.String()
}
