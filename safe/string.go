package safe

import "fmt"

type constantString string

type String interface {
	fmt.Stringer
	Check(required TrustLevel) bool
}

func Const(s constantString) String {
	return s
}

func (constantString) Check(TrustLevel) bool {
	return true
}

func (s constantString) String() string {
	return string(s)
}

type UntrustedString string

func (r UntrustedString) String() string {
	return string(r)
}

func (UntrustedString) Check(required TrustLevel) bool {
	switch required {
	case Untrusted:
		return true
	default:
		return false
	}
}

type URL struct {
	s string
}

func (u URL) String() string {
	return u.s
}

func (URL) Check(required TrustLevel) bool {
	switch required {
	case URLSafe, Untrusted:
		return true
	default:
		return false
	}
}

func EscapeURL(s string) URL {
	u, err := escapeURL(s)
	if err != nil {
		return URL{"about:invalid"}
	}

	return URL{u}
}

type Attribute struct {
	s string
}

func (a Attribute) String() string {
	return a.s
}

func (Attribute) Check(required TrustLevel) bool {
	switch required {
	case AttributeSafe, Untrusted:
		return true
	default:
		return false
	}
}

func EscapeAttribute(s string) Attribute {
	u, err := escapeAttribute(s)
	if err != nil {
		return Attribute{"!invalid"}
	}

	return Attribute{u}
}

type HTML struct {
	s string
}

func (h HTML) String() string {
	return h.s
}

func (HTML) Check(required TrustLevel) bool {
	switch required {
	case HTMLSafe, TextSafe, Untrusted:
		return true
	default:
		return false
	}
}

func EscapeHTML(s string) HTML {
	h, err := escapeHTML(s)
	if err != nil {
		return HTML{"<b>invalid html</b>"}
	}

	return HTML{h}
}

type Text struct {
	s string
}

func (t Text) String() string {
	return t.s
}

func (Text) Check(required TrustLevel) bool {
	switch required {
	case TextSafe, Untrusted:
		return true
	default:
		return false
	}
}

func EscapeText(s string) Text {
	t, err := escapeText(s)
	if err != nil {
		return Text{"invalid text"}
	}

	return Text{t}
}
