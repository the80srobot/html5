package html

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTextNode(t *testing.T) {
	for _, tc := range []struct {
		comment string
		input   *TextNode
		opts    *CompileOptions
		values  []ValueArg
		depth   int
		output  string
	}{
		{
			comment: "one line static",
			input:   &TextNode{Value: FullyTrustedString("Hello, World!")},
			opts:    &CompileOptions{},
			output:  "Hello, World!",
		},
		{
			comment: "two lines static inline",
			input:   &TextNode{Value: FullyTrustedString("Hello, World!"), Width: 1},
			opts:    &CompileOptions{},
			output:  "Hello,\nWorld!",
		},
		{
			comment: "block one line",
			input:   &TextNode{Value: FullyTrustedString("Hello, World!")},
			depth:   1,
			opts:    &CompileOptions{Indent: "  "},
			output:  "Hello, World!",
		},
		{
			comment: "block two lines",
			input:   &TextNode{Value: FullyTrustedString("Hello, World!"), Width: 1},
			depth:   1,
			opts:    &CompileOptions{Indent: "  "},
			output:  "Hello,\n  World!",
		},
		{
			comment: "binding",
			input:   &TextNode{Value: Binding("hello")},
			opts:    &CompileOptions{},
			values:  []ValueArg{{Name: "hello", SafeString: TrustedString("Hello, World!", TextSafe)}},
			output:  "Hello, World!",
		},
		{
			comment: "untrusted binding",
			input:   &TextNode{Value: Binding("hello")},
			opts:    &CompileOptions{},
			values:  []ValueArg{{Name: "hello", SafeString: UntrustedString("<p>Hello, World!</p>")}},
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
