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
				Attributes:  []Attribute{{Name: "id", Constant: FullyTrustedString("hello")}},
				Contents: []Node{
					&TextNode{Constant: FullyTrustedString("Hello, World!"), IndentStyle: Block, Width: 1},
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
					{Name: "href", Constant: FullyTrustedString("#title_1")},
					{Name: "rel", Constant: FullyTrustedString("nofollow")},
					{Name: "target", Constant: FullyTrustedString("_blank")},
				},
				Contents: []Node{&TextNode{Constant: FullyTrustedString("Hello!")}},
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
					{Name: "href", Binding: "href"},
					{Name: "rel", Binding: "rel"},
					{Name: "target", Binding: "target"},
				},
				Contents: []Node{
					&TextNode{StringName: "hello"},
				},
			},
			values: []ValueArg{
				{Name: "href", SafeString: TrustedString("#title_1", URLSafe)},
				{Name: "rel", SafeString: TrustedString("nofollow", FullyTrusted)},
				{Name: "target", SafeString: TrustedString("_blank", AttributeSafe)},
				{Name: "hello", SafeString: TrustedString("Hello!", TextSafe)},
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
