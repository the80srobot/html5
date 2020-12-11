package html

import (
	"bytes"
	"fmt"
	"strings"
)

type templateCompiler struct {
	pending        *bytes.Buffer
	chunks         []chunk
	separateChunks bool
	bindings       BindingSet
}

func (tc *templateCompiler) flush() {
	if tc.pending == nil {
		return
	}
	tc.chunks = append(tc.chunks, staticChunk{tc.pending.String()})
	tc.pending = nil
}

func (tc *templateCompiler) appendChunk(c chunk) {
	tc.flush()
	tc.chunks = append(tc.chunks, c)
}

func (tc *templateCompiler) appendStringBinding(name string, trust StringTrust) {
	tag := tc.bindings.AddString(name, trust)
	tc.appendChunk(stringBindingChunk{stringTag: tag})
}

func (tc *templateCompiler) Write(p []byte) (int, error) {
	if tc.separateChunks {
		tc.flush()
	}
	if tc.pending == nil {
		tc.pending = &bytes.Buffer{}
	}
	return tc.pending.Write(p)
}

func (tc *templateCompiler) WriteString(s string) (int, error) {
	if tc.separateChunks {
		tc.flush()
	}
	if tc.pending == nil {
		tc.pending = bytes.NewBufferString(s)
		return len(s), nil
	}
	return tc.pending.WriteString(s)
}

type tagStyle int16

const (
	tagOpen tagStyle = iota
	tagClose
	tagSelfClose
)

func appendTag(tc *templateCompiler, name string, style tagStyle, attributes ...Attribute) error {
	if style == tagClose {
		_, err := fmt.Fprintf(tc, "</%s>", name)
		return err
	}

	if _, err := fmt.Fprintf(tc, "<%s", name); err != nil {
		return err
	}

	for _, a := range attributes {
		if err := appendAttribute(tc, &a); err != nil {
			return err
		}
	}

	if style == tagSelfClose {
		_, err := fmt.Fprint(tc, "/>")
		return err
	}

	_, err := fmt.Fprint(tc, ">")
	return err
}

func appendAttribute(tc *templateCompiler, a *Attribute) error {
	if _, err := fmt.Fprintf(tc, " %s=\"", a.Name); err != nil {
		return err
	}

	if a.StringName != "" {
		reqTrust, ok := requiredTrustPerAttribute[a.Name]
		if !ok {
			reqTrust = FullyTrusted
		}
		tc.appendStringBinding(a.StringName, reqTrust)
		_, err := fmt.Fprint(tc, "\"")
		return err
	}

	_, err := fmt.Fprintf(tc, "%s\"", a.Constant)
	return err
}

func appendIndent(tc *templateCompiler, depth int, indent string) error {
	for i := 0; i < depth; i++ {
		if _, err := tc.WriteString(indent); err != nil {
			return err
		}
	}
	return nil
}

func appendText(tc *templateCompiler, depth int, text *TextNode, is IndentStyle, indent string) error {
	if text.StringName != "" {
		tag := tc.bindings.AddString(text.StringName, TextSafe)
		tc.appendChunk(textBindingChunk{TextNode: *text, depth: depth, indent: indent, indentStyle: is, stringTag: tag})
		return nil
	}

	return fprintBlockText(tc, depth, text.Width, indent, is, strings.NewReader(text.Constant))
}
