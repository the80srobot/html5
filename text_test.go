package html5

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/the80srobot/html5/bindings"
	"github.com/the80srobot/html5/safe"
)

func TestTextNode(t *testing.T) {
	for _, tc := range []struct {
		comment string
		input   *TextNode
		opts    *CompileOptions
		values  []bindings.BindArg
		output  string
	}{
		{
			comment: "one line static",
			input:   &TextNode{Value: safe.Const("Hello, World!")},
			opts:    &CompileOptions{},
			output:  "Hello, World!",
		},
		{
			comment: "two lines static inline",
			input:   &TextNode{Value: safe.Const("Hello, World!")},
			opts:    &CompileOptions{TextWidth: 1},
			output:  "Hello,\nWorld!",
		},
		{
			comment: "block one line",
			input:   &TextNode{Value: safe.Const("Hello, World!")},
			opts:    &CompileOptions{Indent: "  ", RootDepth: 1},
			output:  "Hello, World!",
		},
		{
			comment: "block two lines",
			input:   &TextNode{Value: safe.Const("Hello, World!")},
			opts:    &CompileOptions{Indent: "  ", TextWidth: 1, RootDepth: 1},
			output:  "Hello,\n  World!",
		},
		{
			comment: "binding",
			input:   &TextNode{Value: bindings.Declare("hello", safe.Default)},
			opts:    &CompileOptions{},
			values:  []bindings.BindArg{{Name: "hello", Value: safe.Const("Hello, World!")}},
			output:  "Hello, World!",
		},
		{
			comment: "untrusted binding",
			input:   &TextNode{Value: bindings.Declare("hello", safe.Default)},
			opts:    &CompileOptions{},
			values:  []bindings.BindArg{{Name: "hello", Value: safe.EscapeText("<p>Hello, World!</p>")}},
			output:  "&lt;p&gt;Hello, World!&lt;/p&gt;",
		},
	} {
		t.Run(tc.comment, func(t *testing.T) {
			if diff := cmp.Diff(tc.output, mustGenerateHTML(t, tc.input, tc.opts, tc.values)); diff != "" {
				t.Errorf("GenerateHTML(%v, %v, %v)\n => (-)wanted vs (+)got:\n%s", tc.input, tc.opts, tc.values, diff)
			}
		})
	}
}
