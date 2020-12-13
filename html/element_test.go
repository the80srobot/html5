package html

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
				Attributes:  []Attribute{{Name: "id", Value: FullyTrustedString("hello")}},
				Contents: []Node{
					&TextNode{Value: FullyTrustedString("Hello, World!"), Width: 1},
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
							&TextNode{Value: FullyTrustedString("Span")},
						},
					},
					&ElementNode{
						Name:        "span",
						IndentStyle: Block,
						Contents: []Node{
							&TextNode{Value: FullyTrustedString("Span")},
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
					{Name: "href", Value: FullyTrustedString("#title_1")},
					{Name: "rel", Value: FullyTrustedString("nofollow")},
					{Name: "target", Value: FullyTrustedString("_blank")},
				},
				Contents: []Node{&TextNode{Value: FullyTrustedString("Hello!")}},
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
					{Name: "href", Value: Binding("href")},
					{Name: "rel", Value: Binding("rel")},
					{Name: "target", Value: Binding("target")},
				},
				Contents: []Node{
					&TextNode{Value: Binding("hello")},
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
