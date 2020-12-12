package html5

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/the80srobot/html5/html"
)

func TestElements(t *testing.T) {
	for _, tc := range []struct {
		comment string
		opts    *html.CompileOptions
		input   html.Node
		values  []html.ValueArg
		output  string
	}{
		{
			comment: "html page",
			opts:    &html.Tidy,
			input: HTML(
				Head(),
				Element("body", Indent(html.Block),
					Element("h1", Text(html.FullyTrustedString("Hello, "), html.Bind("user_name"))))),
			values: []html.ValueArg{{Name: "user_name", SafeString: html.UntrustedString("Bob")}},
			output: `<!doctype html>
<html>
  <head>
    <meta charset="utf-8">
  </head><body>
    <h1>Hello, Bob</h1>
  </body>
</html>`, // TODO: newline after head would be nicer.
		},
	} {
		t.Run(tc.comment, func(t *testing.T) {
			doc, err := Compile(tc.input, tc.opts)
			if err != nil {
				t.Fatalf("Compile(%v, %v): %v", tc.input, tc.opts, err)
			}

			var sb strings.Builder
			if err := doc.GenerateHTML(&sb, tc.values...); err != nil {
				t.Fatalf("GenerateHTML(): %v", err)
			}

			if diff := cmp.Diff(tc.output, sb.String()); diff != "" {
				t.Errorf("GenerateHTML(%v, %v, %v)\n => (-)wanted vs (+)got:\n%s", tc.input, tc.opts, tc.values, diff)
			}
		})
	}
}
