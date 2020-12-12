package html

type ElementNode struct {
	Name                string
	Attributes          []Attribute
	Contents            []Node
	IndentStyle         IndentStyle
	XMLStyleSelfClosing bool
}

func (e *ElementNode) deduplicateAttributes() {
	var attrs []Attribute
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

func (e *ElementNode) compile(db *templateCompiler, depth int, opts *CompileOptions) error {
	e.deduplicateAttributes()

	if e.IndentStyle == Block && !opts.Compact {
		if err := fprintRawNewline(db, depth, opts.Indent); err != nil {
			return err
		}
		depth++
	}

	// Handle the xHTML style (produce <br/> instead of <br>).
	openingTag := tagOpen
	if len(e.Contents) == 0 && e.XMLStyleSelfClosing {
		openingTag = tagSelfClose
	}

	if err := appendTag(db, e.Name, openingTag, e.Attributes...); err != nil {
		return err
	}
	if len(e.Contents) == 0 {
		return nil
	}

	for _, c := range e.Contents {
		if err := c.compile(db, depth, opts); err != nil {
			return err
		}
	}

	if e.IndentStyle == Block && !opts.Compact {
		depth--
		if err := fprintRawNewline(db, depth, opts.Indent); err != nil {
			return err
		}
	}

	return appendTag(db, e.Name, tagClose)
}
