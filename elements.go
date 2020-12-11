package html5

import "github.com/the80srobot/html5/html"

type Document struct {
	*html.Template
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

func ConstantAttribute(name, constant string) *html.Attribute {
	return &html.Attribute{Name: name, Constant: constant}
}

func HTML(head, body html.Node, opts ...Option) *html.MultiNode {
	e := &html.ElementNode{
		Name:        "html",
		Contents:    []html.Node{head, body},
		IndentStyle: html.Block,
	}
	applyOptions(e, opts...)
	return &html.MultiNode{
		Contents: []html.Node{
			&html.TextNode{Constant: "<!doctype html>", Width: 0, IndentStyle: html.Block},
			e,
		},
	}
}

func Head(contents []html.Node, opts ...Option) *html.ElementNode {
	e := &html.ElementNode{
		Name:        "head",
		IndentStyle: html.Block,
		Contents: []html.Node{
			Meta(&html.Attribute{Name: "charset", Constant: "utf-8"}),
		},
	}
	e.Contents = append(e.Contents, contents...)
	applyOptions(e)
	return e
}

func Meta(opts ...Option) *html.ElementNode {
	e := &html.ElementNode{
		Name:        "meta",
		IndentStyle: html.Block,
	}
	applyOptions(e)
	return e
}

func P(contents []html.Node, opts ...Option) *html.ElementNode {
	e := &html.ElementNode{
		Name:     "p",
		Contents: contents,
	}
	applyOptions(e, opts...)
	return e
}

func Constant(constant string, opts ...Option) *html.TextNode {
	e := &html.TextNode{Constant: constant}
	applyOptions(e)
	return e
}

func Text(name string, opts ...Option) *html.TextNode {
	e := &html.TextNode{StringName: name}
	applyOptions(e)
	return e
}
