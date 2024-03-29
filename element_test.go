package html5

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/the80srobot/html5/bindings"
	"github.com/the80srobot/html5/safe"
)

func TestElementNode(t *testing.T) {
	opts := Tidy
	opts.TextWidth = 1 // Force line breaks after every word.

	for _, tc := range []struct {
		comment string
		input   *ElementNode
		opts    *CompileOptions
		values  []bindings.BindArg
		output  string
	}{
		{
			comment: "paragraph",
			input: Element("p",
				Attribute("id", safe.Const("hello")),
				Text(safe.Const("Hello, World!"))),
			opts: &opts,
			output: `<p id="hello">
  Hello,
  World!
</p>`,
		},
		{
			comment: "block indent",
			input: Element("p",
				Element("span", Indent(Block), Text(safe.Const("Span"))),
				Element("span", Indent(Block), Text(safe.Const("Span")))),
			opts: &opts,
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
			input: Element("a",
				Attribute("href", safe.Const("#title_1")),
				Attribute("rel", safe.Const("nofollow")),
				Attribute("target", safe.Const("_blank")),
				Text(safe.Const("Hello!"))),
			opts:   &opts,
			output: "<a href=\"#title_1\" rel=\"nofollow\" target=\"_blank\">Hello!</a>",
		},
		{
			comment: "bindings",
			input: Element("a",
				Attribute("href", bindings.Declare("href", safe.Default)),
				Attribute("rel", bindings.Declare("rel", safe.Default)),
				Attribute("target", bindings.Declare("target", safe.Default)),
				Text(bindings.Declare("hello", safe.Default)),
			),
			values: []bindings.BindArg{
				{Name: "href", Value: safe.Const("#title_1")},
				{Name: "rel", Value: safe.Const("nofollow")},
				{Name: "target", Value: safe.Const("_blank")},
				{Name: "hello", Value: safe.Const("Hello!")},
			},
			opts:   &opts,
			output: "<a href=\"#title_1\" rel=\"nofollow\" target=\"_blank\">Hello!</a>",
		},
	} {
		opt := cmpopts.AcyclicTransformer("multiline", func(s string) []string {
			return strings.Split(s, "\n")
		})
		t.Run(tc.comment, func(t *testing.T) {
			if diff := cmp.Diff(tc.output, mustGenerateHTML(t, tc.input, tc.opts, tc.values), opt); diff != "" {
				t.Errorf("GenerateHTML(%v, %v, %v)\n => (-)wanted vs (+)got:\n%s", tc.input, tc.opts, tc.values, diff)
			}
		})
	}
}
