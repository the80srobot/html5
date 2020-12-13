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

func (tc *templateCompiler) freshLine() bool {
	if tc.pending == nil {
		return len(tc.chunks) == 0
	}

	p := tc.pending.Bytes()
	for i := len(p) - 1; i >= 0; i-- {
		switch p[i] {
		case '\n':
			return true
		case ' ', '\t':
			continue
		default:
			return false
		}
	}

	return false
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
	tag := tc.bindings.DeclareString(name, trust)
	tc.appendChunk(stringBindingChunk{stringTag: tag})
}

func (tc *templateCompiler) Write(p []byte) (int, error) {
	tc.ensureBuffer()
	return tc.pending.Write(p)
}

func (tc *templateCompiler) WriteString(s string) (int, error) {
	tc.ensureBuffer()
	return tc.pending.WriteString(s)
}

func (tc *templateCompiler) ensureBuffer() {
	if tc.separateChunks {
		tc.flush()
	}
	if tc.pending == nil {
		tc.pending = &bytes.Buffer{}
	}
}

func (tc *templateCompiler) appendNewLine(depth int, indent string) error {
	tc.ensureBuffer()
	if _, err := tc.pending.Write([]byte{'\n'}); err != nil {
		return err
	}
	for i := 0; i < depth; i++ {
		if _, err := tc.pending.WriteString(indent); err != nil {
			return err
		}
	}
	return nil
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

	// Different attributes require different levels of trust (e.g. href
	// contains URLs).
	reqTrust, ok := requiredTrustPerAttribute[a.Name]
	if !ok {
		reqTrust = FullyTrusted
	}

	if a.Value.Constant() {
		constant, err := a.Value.Convert(reqTrust)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(tc, "%s\"", constant)
		return err
	}

	tc.appendStringBinding(a.Value.binding, reqTrust)
	_, err := fmt.Fprint(tc, "\"")
	return err
}

func appendText(tc *templateCompiler, depth int, text *TextNode, indent string) error {
	if text.Value.Constant() {
		constant, err := text.Value.Convert(TextSafe)
		if err != nil {
			return err
		}
		return fprintBlockText(tc, depth, text.Width, indent, strings.NewReader(constant))
	}

	tag := tc.bindings.DeclareString(text.Value.binding, TextSafe)
	tc.appendChunk(textBindingChunk{TextNode: *text, depth: depth, indent: indent, stringTag: tag})
	return nil
}
