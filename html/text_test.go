package html

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
		depth   int
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
			input:   &TextNode{Value: safe.Const("Hello, World!"), Width: 1},
			opts:    &CompileOptions{},
			output:  "Hello,\nWorld!",
		},
		{
			comment: "block one line",
			input:   &TextNode{Value: safe.Const("Hello, World!")},
			depth:   1,
			opts:    &CompileOptions{Indent: "  "},
			output:  "Hello, World!",
		},
		{
			comment: "block two lines",
			input:   &TextNode{Value: safe.Const("Hello, World!"), Width: 1},
			depth:   1,
			opts:    &CompileOptions{Indent: "  "},
			output:  "Hello,\n  World!",
		},
		{
			comment: "binding",
			input:   &TextNode{Value: bindings.Declare("hello")},
			opts:    &CompileOptions{},
			values:  []bindings.BindArg{{Name: "hello", Value: safe.Const("Hello, World!")}},
			output:  "Hello, World!",
		},
		{
			comment: "untrusted binding",
			input:   &TextNode{Value: bindings.Declare("hello")},
			opts:    &CompileOptions{},
			values:  []bindings.BindArg{{Name: "hello", Value: safe.EscapeText("<p>Hello, World!</p>")}},
			output:  "&lt;p&gt;Hello, World!&lt;/p&gt;",
		},
	} {
		t.Run(tc.comment, func(t *testing.T) {
			if diff := cmp.Diff(tc.output, mustGenerateHTML(t, tc.input, tc.depth, tc.opts, tc.values)); diff != "" {
				t.Errorf("GenerateHTML(%v, %v, %v, %v)\n => (-)wanted vs (+)got:\n%s", tc.input, tc.depth, tc.opts, tc.values, diff)
			}
		})
	}
}
