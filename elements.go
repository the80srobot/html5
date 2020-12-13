package html5

import (
	"io"

	"github.com/the80srobot/html5/html"
)

type Document struct {
	*html.Template
}

func (d *Document) GenerateHTML(w io.Writer, values ...html.ValueArg) error {
	vs, err := d.Bingings.Bind(values...)
	if err != nil {
		return err
	}
	return d.Template.GenerateHTML(w, vs)
}

func Compile(node html.Node, opts *html.CompileOptions) (*Document, error) {
	t, err := html.Compile(node, 0, opts)
	if err != nil {
		return nil, err
	}
	return &Document{t}, nil
}

func MustCompile(node html.Node, opts *html.CompileOptions) *Document {
	d, err := Compile(node, opts)
	if err != nil {
		panic(err)
	}
	return d
}

func HTML(values ...interface{}) *html.MultiNode {
	values = append(values, Indent(html.Block))
	e := Element("html", values...)
	return &html.MultiNode{
		Contents: []html.Node{
			&html.RawNode{HTML: html.FullyTrustedString("<!doctype html>\n")},
			e,
		},
	}
}

func Head(values ...interface{}) *html.ElementNode {
	values = append(values, Indent(html.Block), Meta(Attribute("charset", html.FullyTrustedString("utf-8"))))
	return Element("head", values...)
}

func Meta(opts ...Option) *html.ElementNode {
	e := &html.ElementNode{
		Name:        "meta",
		IndentStyle: html.Block,
		SelfClosing: true,
	}
	applyOptions(e, opts...)
	return e
}

func Element(name string, values ...interface{}) *html.ElementNode {
	e := &html.ElementNode{Name: name}
	for _, value := range values {
		switch value := value.(type) {
		case Option:
			value.Apply(e)
		case html.Node:
			e.Contents = append(e.Contents, value)
		default:
			panic("invalid argument (must be option or html.Node)")
		}
	}
	return e
}

func Attribute(name string, value html.SafeString) *html.Attribute {
	return &html.Attribute{Name: name, Value: value}
}

func Text(parts ...html.SafeString) *html.MultiNode {
	e := &html.MultiNode{Contents: make([]html.Node, 0, len(parts))}
	for _, part := range parts {
		e.Contents = append(e.Contents, &html.TextNode{Value: part})
	}
	return e
}
