package html

import (
	"strings"
	"testing"
)

func mustGenerateHTML(t testing.TB, n Node, depth int, opts *CompileOptions, values []ValueArg) string {
	t.Helper()
	var sb strings.Builder
	var bs BindingSet
	tmpl, err := Compile(n, depth, &bs, opts)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}

	t.Logf("compiled template: %v", tmpl)

	vs, err := tmpl.Bingings.Bind(values...)
	if err != nil {
		t.Fatalf("Bind: %v", err)
	}

	t.Logf("bound values: %v", vs)

	if err := tmpl.GenerateHTML(&sb, vs); err != nil {
		t.Fatalf("GenerateHTML: %v", err)
	}
	return sb.String()
}
