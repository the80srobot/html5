package html5

import (
	"fmt"
	"io"

	"github.com/the80srobot/html5/html"
)

type Document struct {
	*html.Template
}

// An Input is an Option or an html.Node.
type Input interface{}

// A String is an html.SafeString or a plain string, defaulting to
// untrustred.
type String interface{}

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

func HTML(values ...Input) *html.MultiNode {
	values = append(values, Indent(html.Block))
	e := Element("html", values...)
	return &html.MultiNode{
		Contents: []html.Node{
			&html.RawNode{HTML: html.FullyTrustedString("<!doctype html>\n")},
			e,
		},
	}
}

func Head(values ...Input) *html.ElementNode {
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

func Element(name string, values ...Input) *html.ElementNode {
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

func Attribute(name string, value String) *html.Attribute {
	return &html.Attribute{Name: name, Value: wrap(value)}
}

func Text(parts ...String) *html.MultiNode {
	e := &html.MultiNode{Contents: make([]html.Node, 0, len(parts))}
	for _, part := range parts {
		e.Contents = append(e.Contents, &html.TextNode{Value: wrap(part)})
	}
	return e
}

func Multi(nodes ...html.Node) *html.MultiNode {
	return &html.MultiNode{Contents: nodes}
}

func Raw(s String) *html.RawNode {
	return &html.RawNode{HTML: wrap(s)}
}

func wrap(s String) html.SafeString {
	switch s := s.(type) {
	case html.SafeString:
		return s
	case string:
		return html.UntrustedString(s)
	default:
		panic(fmt.Sprintf("String must be either string or html.SafeString, not %v", s))
	}
}

func Safe(s string) html.SafeString {
	return html.FullyTrustedString(s)
}
