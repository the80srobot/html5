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
				Attribute("lang", html.FullyTrustedString("en")),
				Head(
					Meta(Attribute("author", html.Binding("author"))),
					Meta(Attribute("description", html.Binding("description")))),
				Element("body", Indent(html.Block),
					Element("h1", Text(html.FullyTrustedString("Hello, "), html.Binding("user_name"))),
					Element("p", Indent(html.Block), Text(html.FullyTrustedString("Welcome to our website."))))),
			values: []html.ValueArg{
				{Name: "author", SafeString: html.FullyTrustedString("Nakatomi Corporation")},
				{Name: "description", SafeString: html.FullyTrustedString("Our glorious website")},
				{Name: "user_name", SafeString: html.UntrustedString("Bob")},
			},
			output: `<!doctype html>
<html lang="en">
  <head>
    <meta author="Nakatomi Corporation">
    <meta description="Our glorious website">
    <meta charset="utf-8">
  </head>
  <body>
    <h1>Hello, Bob</h1>
    <p>
      Welcome to our website.
    </p>
  </body>
</html>`,
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
