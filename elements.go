package html5

import (
	"io"

	"github.com/the80srobot/html5/bindings"
	"github.com/the80srobot/html5/html"
	"github.com/the80srobot/html5/safe"
)

// An Input is an Option or an html.Node.
type Input interface{}

// A String is an html.SafeString or a plain string, defaulting to
// untrustred.
type String interface{}

func GenerateHTML(w io.Writer, t *html.Template, values ...bindings.BindArg) error {
	vm, err := t.Bindings.Bind()
	if err != nil {
		return err
	}
	bindings.Bind(vm, values...)
	return t.GenerateHTML(w, vm)
}

func Compile(node html.Node, m *bindings.Map, opts *html.CompileOptions) (*html.Template, error) {
	t, err := html.Compile(node, 0, m, opts)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func MustCompile(node html.Node, m *bindings.Map, opts *html.CompileOptions) *html.Template {
	d, err := Compile(node, m, opts)
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
			&html.RawNode{HTML: safe.Const("<!doctype html>\n")},
			e,
		},
	}
}

func Head(values ...Input) *html.ElementNode {
	values = append(values, Indent(html.Block), Meta(Attribute("charset", safe.Const("utf-8"))))
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

func Attribute(name string, value safe.String) *html.Attribute {
	return &html.Attribute{Name: name, Value: value}
}

func Text(parts ...safe.String) *html.MultiNode {
	e := &html.MultiNode{Contents: make([]html.Node, 0, len(parts))}
	for _, part := range parts {
		e.Contents = append(e.Contents, &html.TextNode{Value: part})
	}
	return e
}

func Multi(nodes ...html.Node) *html.MultiNode {
	return &html.MultiNode{Contents: nodes}
}

func Raw(s safe.String) *html.RawNode {
	return &html.RawNode{HTML: s}
}
