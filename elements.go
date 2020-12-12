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

func ConstantAttribute(name string, constant html.SafeString) *html.Attribute {
	return &html.Attribute{Name: name, Value: constant}
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
			&html.TextNode{Value: html.FullyTrustedString("<!doctype html>"), Width: 0, IndentStyle: html.Block},
			e,
		},
	}
}

func Head(contents []html.Node, opts ...Option) *html.ElementNode {
	e := &html.ElementNode{
		Name:        "head",
		IndentStyle: html.Block,
		Contents: []html.Node{
			Meta(&html.Attribute{Name: "charset", Value: html.FullyTrustedString("utf-8")}),
		},
	}
	e.Contents = append(e.Contents, contents...)
	applyOptions(e)
	return e
}

func Body(contents []html.Node, opts ...Option) *html.ElementNode {
	e := &html.ElementNode{
		Name:        "body",
		IndentStyle: html.Block,
		Contents:    contents,
	}
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
