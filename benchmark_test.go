package html5

import (
	"html/template"
	"io/ioutil"
	"testing"

	safetemplate "github.com/google/safehtml/template"
	"github.com/the80srobot/html5/html"
)

const smallTemplate = `<!doctype html>
<html lang="en">
  <head>
  <meta author="{{.Author}}">
  <meta description="{{.Description}}">
  <meta charset="utf-8">
  </head>
<body>
  <h1>Hello, {{.User}}</h1>
  <p>
    Welcome to our website.
  </p>
</body>
</html>`

type smallTemplateData struct {
	Author      string
	Description string
	User        string
}

func BenchmarkSmallTemplate(b *testing.B) {
	tmpl := template.Must(template.New("test").Parse(smallTemplate))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := tmpl.Execute(ioutil.Discard, smallTemplateData{"Nakatomi Corporation", "Our glorious website", "Bob"}); err != nil {
			b.Fatalf("Execute failed: %v", err)
		}
	}
}

func BenchmarkSmallPage(b *testing.B) {
	page := HTML(
		Attribute("lang", html.FullyTrustedString("en")),
		Head(
			Meta(Attribute("author", html.Binding("author"))),
			Meta(Attribute("description", html.Binding("description")))),
		Element("body", Indent(html.Block),
			Element("h1", Text(html.FullyTrustedString("Hello, "), html.Binding("user_name"))),
			Element("p", Indent(html.Block), Text(html.FullyTrustedString("Welcome to our website.")))))
	doc, err := Compile(page, &html.Compact)

	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := doc.GenerateHTML(ioutil.Discard, []html.ValueArg{
			{Name: "author", SafeString: html.FullyTrustedString("Nakatomi Corporation")},
			{Name: "description", SafeString: html.FullyTrustedString("Our glorious website")},
			{Name: "user_name", SafeString: html.UntrustedString("Bob")},
		}...)
		if err != nil {
			b.Fatalf("GenerateHTML failed: %v", err)
		}
	}
}

type smallSafeTemplateData struct {
	Author      safetemplate.TrustedSource
	Description safetemplate.TrustedSource
	User        string
}

// The safehtml benchmark is disabled, because it seems like the safehtml
// package is not capable of generating the same HTML as html/template and
// html5. There seems to be no way to force safehtml to interpolate strings in
// attributes, no matter how well-attested their origin is.

// func BenchmarkSmallSafeHTMLTemplate(b *testing.B) {
// 	tmpl := safetemplate.Must(safetemplate.New("test").Parse(smallTemplate))
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		if err := tmpl.Execute(ioutil.Discard, smallSafeTemplateData{
// 			safetemplate.TrustedSourceFromConstant("Nakatomi Corporation"),
// 			safetemplate.TrustedSourceFromConstant("Our glorious website"),
// 			"Bob"}); err != nil {
// 			b.Fatalf("Execute failed: %v", err)
// 		}
// 	}
// }
