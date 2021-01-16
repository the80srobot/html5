package html5

import (
	"strings"
	"testing"

	"github.com/the80srobot/html5/bindings"
	"github.com/the80srobot/html5/safe"
)

func mustGenerateHTML(t testing.TB, n Node, depth int, opts *CompileOptions, values []bindings.BindArg) string {
	t.Helper()
	var sb strings.Builder
	var m bindings.Map
	tmpl, err := Compile(n, depth, &m, opts)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}

	t.Logf("Compiled template: %v", tmpl)

	vm, err := tmpl.Bindings.Bind()
	if err != nil {
		t.Fatalf("Bind: %v", err)
	}
	if err := bindings.Bind(vm, values...); err != nil {
		t.Fatalf("Bind: %v", err)
	}

	t.Logf("Bound values: %v", vm)

	if err := tmpl.GenerateHTML(&sb, vm); err != nil {
		t.Fatalf("GenerateHTML: %v", err)
	}
	return sb.String()
}

type valueArg struct {
	Name        string
	Value       safe.String
	Subsections [][]valueArg
}
