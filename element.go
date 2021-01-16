package html5

import (
	"errors"
	"fmt"
)

// ElementNode represents an HTML element, like <p>.
type ElementNode struct {
	Name                string
	Attributes          []AttributeNode
	Contents            []Node
	IndentStyle         IndentStyle
	SelfClosing         bool
	XMLStyleSelfClosing bool
}

func Element(name string, contents ...Content) *ElementNode {
	e, ok := elementPrototypes[name]
	if !ok {
		e = ElementNode{Name: name}
	}
	for _, c := range contents {
		c.Apply(&e)
	}
	return &e
}

func (e *ElementNode) Apply(n Node) error {
	switch n := n.(type) {
	case *ElementNode:
		n.Contents = append(n.Contents, e)
	case *MultiNode:
		n.Contents = append(n.Contents, e)
	default:
		return fmt.Errorf("ElementNode can only be applied to ElementNode or MultiNode, got %v", n)
	}
	return nil
}

func (e *ElementNode) deduplicateAttributes() {
	var attrs []AttributeNode
	seen := make(map[string]struct{}, len(e.Attributes))
	for i := len(e.Attributes) - 1; i >= 0; i-- {
		if _, ok := seen[e.Attributes[i].Name]; ok {
			continue
		}
		seen[e.Attributes[i].Name] = struct{}{}
		attrs = append(attrs, e.Attributes[i])
	}

	// attrs are deduplicated and also reversed.
	l := len(attrs)
	for i := 0; i < l/2; i++ {
		attrs[i], attrs[l-i-1] = attrs[l-i-1], attrs[i]
	}
	e.Attributes = attrs
}

func (e *ElementNode) compile(tc *templateCompiler, depth int, opts *CompileOptions) error {
	e.deduplicateAttributes()

	isBlock := e.IndentStyle == Block && !opts.Compact

	// Block elements always start on a new line.
	if isBlock && !tc.freshLine() {
		if err := tc.appendNewLine(depth, opts.Indent); err != nil {
			return err
		}
	}

	// Handle the xHTML style (produce <br /> instead of <br>).
	openingTag := tagOpen
	if e.XMLStyleSelfClosing {
		openingTag = tagSelfClose
	}

	if err := appendTag(tc, e.Name, openingTag, e.Attributes...); err != nil {
		return err
	}

	if e.SelfClosing || e.XMLStyleSelfClosing {
		if len(e.Contents) != 0 {
			return errors.New("self-closing element cannot have contents")
		}
		return nil
	}

	if isBlock {
		depth++
		// Block element contents are indented by an additional level.
		if err := tc.appendNewLine(depth, opts.Indent); err != nil {
			return err
		}
	}

	for _, c := range e.Contents {
		if err := c.compile(tc, depth, opts); err != nil {
			return err
		}
	}

	if isBlock {
		depth--
		if err := tc.appendNewLine(depth, opts.Indent); err != nil {
			return err
		}
	}

	if err := appendTag(tc, e.Name, tagClose); err != nil {
		return err
	}

	return nil
}

var elementPrototypes = map[string]ElementNode{
	// HTML5 void elements: self-closing by default.
	"area":    {Name: "area", IndentStyle: Block, SelfClosing: true},
	"base":    {Name: "base", IndentStyle: Inline, SelfClosing: true},
	"br":      {Name: "br", IndentStyle: Inline, SelfClosing: true},
	"command": {Name: "command", IndentStyle: Inline, SelfClosing: true},
	"embed":   {Name: "embed", IndentStyle: Block, SelfClosing: true},
	"hr":      {Name: "hr", IndentStyle: Block, SelfClosing: true},
	"img":     {Name: "img", IndentStyle: Inline, SelfClosing: true},
	"input":   {Name: "input", IndentStyle: Block, SelfClosing: true},
	"keygen":  {Name: "keygen", IndentStyle: Block, SelfClosing: true},
	"link":    {Name: "link", IndentStyle: Block, SelfClosing: true},
	"meta":    {Name: "meta", IndentStyle: Block, SelfClosing: true},
	"param":   {Name: "param", IndentStyle: Block, SelfClosing: true},
	"source":  {Name: "source", IndentStyle: Block, SelfClosing: true},
	"track":   {Name: "track", IndentStyle: Block, SelfClosing: true},
	"wbr":     {Name: "wbr", IndentStyle: Inline, SelfClosing: true},

	// Block indent by default.
	"article":    {Name: "article", IndentStyle: Block, SelfClosing: false},
	"address":    {Name: "address", IndentStyle: Block, SelfClosing: false},
	"aside":      {Name: "aside", IndentStyle: Block, SelfClosing: false},
	"body":       {Name: "body", IndentStyle: Block, SelfClosing: false},
	"blockquote": {Name: "blockquote", IndentStyle: Block, SelfClosing: false},
	"button":     {Name: "button", IndentStyle: Block, SelfClosing: false},
	"canvas":     {Name: "canvas", IndentStyle: Block, SelfClosing: false},
	"caption":    {Name: "caption", IndentStyle: Block, SelfClosing: false},
	"code":       {Name: "code", IndentStyle: Block, SelfClosing: false},
	"colgroup":   {Name: "colgroup", IndentStyle: Block, SelfClosing: false},
	"datalist":   {Name: "datalist", IndentStyle: Block, SelfClosing: false},
	"dl":         {Name: "dl", IndentStyle: Block, SelfClosing: false},
	"fieldset":   {Name: "fieldset", IndentStyle: Block, SelfClosing: false},
	"figcaption": {Name: "figcaption", IndentStyle: Block, SelfClosing: false},
	"figure":     {Name: "figure", IndentStyle: Block, SelfClosing: false},
	"footer":     {Name: "footer", IndentStyle: Block, SelfClosing: false},
	"form":       {Name: "form", IndentStyle: Block, SelfClosing: false},
	"h1":         {Name: "h1", IndentStyle: Block, SelfClosing: false},
	"h2":         {Name: "h2", IndentStyle: Block, SelfClosing: false},
	"h3":         {Name: "h3", IndentStyle: Block, SelfClosing: false},
	"h4":         {Name: "h4", IndentStyle: Block, SelfClosing: false},
	"h5":         {Name: "h5", IndentStyle: Block, SelfClosing: false},
	"h6":         {Name: "h6", IndentStyle: Block, SelfClosing: false},
	"head":       {Name: "head", IndentStyle: Block, SelfClosing: false},
	"header":     {Name: "header", IndentStyle: Block, SelfClosing: false},
	"html":       {Name: "html", IndentStyle: Block, SelfClosing: false},
	"hgroup":     {Name: "hgroup", IndentStyle: Block, SelfClosing: false},
	"iframe":     {Name: "iframe", IndentStyle: Block, SelfClosing: false},
	"legend":     {Name: "legend", IndentStyle: Block, SelfClosing: false},
	"nav":        {Name: "nav", IndentStyle: Block, SelfClosing: false},
	"p":          {Name: "p", IndentStyle: Block, SelfClosing: false},
	"ol":         {Name: "ol", IndentStyle: Block, SelfClosing: false},
	"optgroup":   {Name: "optgroup", IndentStyle: Block, SelfClosing: false},
	"samp":       {Name: "samp", IndentStyle: Block, SelfClosing: false},
	"script":     {Name: "script", IndentStyle: Block, SelfClosing: false},
	"section":    {Name: "section", IndentStyle: Block, SelfClosing: false},
	"style":      {Name: "style", IndentStyle: Block, SelfClosing: false},
	"table":      {Name: "table", IndentStyle: Block, SelfClosing: false},
	"tbody":      {Name: "tbody", IndentStyle: Block, SelfClosing: false},
	"textarea":   {Name: "textarea", IndentStyle: Block, SelfClosing: false},
	"tfoot":      {Name: "tfoot", IndentStyle: Block, SelfClosing: false},
	"thead":      {Name: "thead", IndentStyle: Block, SelfClosing: false},
	"title":      {Name: "title", IndentStyle: Block, SelfClosing: false},
	"ul":         {Name: "ul", IndentStyle: Block, SelfClosing: false},

	// Everything else defaults to paired tags (non-void), inline.
}
