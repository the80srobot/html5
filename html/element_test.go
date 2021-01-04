package html

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/the80srobot/html5/bindings"
	"github.com/the80srobot/html5/safe"
)

func TestElementNode(t *testing.T) {
	for _, tc := range []struct {
		comment string
		input   *ElementNode
		opts    *CompileOptions
		values  []bindings.BindArg
		depth   int
		output  string
	}{
		{
			comment: "paragraph",
			input: &ElementNode{
				IndentStyle: Block,
				Name:        "p",
				Attributes:  []Attribute{{Name: "id", Value: safe.Const("hello")}},
				Contents: []Node{
					&TextNode{Value: safe.Const("Hello, World!"), Width: 1},
				}},
			opts: &Tidy,
			output: `<p id="hello">
  Hello,
  World!
</p>`,
		},
		{
			comment: "block indent",
			input: &ElementNode{
				IndentStyle: Block,
				Name:        "p",
				Contents: []Node{
					&ElementNode{
						Name:        "span",
						IndentStyle: Block,
						Contents: []Node{
							&TextNode{Value: safe.Const("Span")},
						},
					},
					&ElementNode{
						Name:        "span",
						IndentStyle: Block,
						Contents: []Node{
							&TextNode{Value: safe.Const("Span")},
						},
					},
				},
			},
			opts: &Tidy,
			output: `<p>
  <span>
    Span
  </span>
  <span>
    Span
  </span>
</p>`,
		},
		{
			comment: "multiple attributes",
			input: &ElementNode{
				IndentStyle: Inline,
				Name:        "a",
				Attributes: []Attribute{
					{Name: "href", Value: safe.Const("#title_1")},
					{Name: "rel", Value: safe.Const("nofollow")},
					{Name: "target", Value: safe.Const("_blank")},
				},
				Contents: []Node{&TextNode{Value: safe.Const("Hello!")}},
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
					{Name: "href", Value: bindings.Declare("href")},
					{Name: "rel", Value: bindings.Declare("rel")},
					{Name: "target", Value: bindings.Declare("target")},
				},
				Contents: []Node{
					&TextNode{Value: bindings.Declare("hello")},
				},
			},
			values: []bindings.BindArg{
				{Name: "href", Value: safe.Const("#title_1")},
				{Name: "rel", Value: safe.Const("nofollow")},
				{Name: "target", Value: safe.Const("_blank")},
				{Name: "hello", Value: safe.Const("Hello!")},
			},
			opts:   &Tidy,
			output: "<a href=\"#title_1\" rel=\"nofollow\" target=\"_blank\">Hello!</a>",
		},
	} {
		opt := cmpopts.AcyclicTransformer("multiline", func(s string) []string {
			return strings.Split(s, "\n")
		})
		t.Run(tc.comment, func(t *testing.T) {
			if diff := cmp.Diff(tc.output, mustGenerateHTML(t, tc.input, tc.depth, tc.opts, tc.values), opt); diff != "" {
				t.Errorf("GenerateHTML(%v, %v, %v, %v)\n => (-)wanted vs (+)got:\n%s", tc.input, tc.depth, tc.opts, tc.values, diff)
			}
		})
	}
}
