package html

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestElementNode(t *testing.T) {
	for _, tc := range []struct {
		comment string
		input   *ElementNode
		opts    *CompileOptions
		values  []ValueArg
		depth   int
		output  string
	}{
		{
			comment: "paragraph",
			input: &ElementNode{
				IndentStyle: Block,
				Name:        "p",
				Attributes:  []Attribute{{Name: "id", Constant: "hello"}},
				Contents: []Node{
					&TextNode{Constant: "Hello, World!", IndentStyle: Block, Width: 1},
				}},
			opts:   &Tidy,
			output: "\n<p id=\"hello\">\n  Hello,\n  World!\n</p>",
		},
		{
			comment: "multiple attributes",
			input: &ElementNode{
				IndentStyle: Inline,
				Name:        "a",
				Attributes: []Attribute{
					{Name: "href", Constant: "#title_1"},
					{Name: "rel", Constant: "nofollow"},
					{Name: "target", Constant: "_blank"},
				},
				Contents: []Node{&TextNode{Constant: "Hello!"}},
			},
			opts:   &Tidy,
			output: "<a href=\"#title_1\" rel=\"nofollow\" target=\"_blank\">Hello!</a>",
		},
		{
			comment: "bindings",
			input: &ElementNode{
				IndentStyle: Inline,
				Name:        "a",
				Attributes: []Attribute{
					{Name: "href", StringName: "href"},
					{Name: "rel", StringName: "rel"},
					{Name: "target", StringName: "target"},
				},
				Contents: []Node{
					&TextNode{StringName: "hello"},
				},
			},
			values: []ValueArg{
				{Name: "href", StringValue: "#title_1", StringTrust: URLSafe},
				{Name: "rel", StringValue: "nofollow", StringTrust: FullyTrusted},
				{Name: "target", StringValue: "_blank", StringTrust: AttributeSafe},
				{Name: "hello", StringValue: "Hello!", StringTrust: TextSafe},
			},
			opts:   &Tidy,
			output: "<a href=\"#title_1\" rel=\"nofollow\" target=\"_blank\">Hello!</a>",
		},
	} {
		t.Run(tc.comment, func(t *testing.T) {
			if diff := cmp.Diff(tc.output, mustGenerateHTML(t, tc.input, tc.depth, tc.opts, tc.values)); diff != "" {
				t.Errorf("GenerateHTML(%v, %v, %v, %v)\n => (-)wanted vs (+)got:\n%s", tc.input, tc.depth, tc.opts, tc.values, diff)
			}
		})
	}
}
