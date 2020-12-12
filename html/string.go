package html

import (
	"errors"
	"fmt"
	"html"

	"github.com/google/safehtml"
)

type StringTrust int16

var ErrStringUntrusted = errors.New("cannnot satisfy trust level")

const (
	Untrusted StringTrust = iota
	HTMLSafe
	TextSafe
	URLSafe
	AttributeSafe
	FullyTrusted
)

func (l StringTrust) String() string {
	switch l {
	case Untrusted:
		return "Untrusted"
	case HTMLSafe:
		return "HTMLSafe"
	case TextSafe:
		return "TextSafe"
	case URLSafe:
		return "URLSafe"
	case AttributeSafe:
		return "AttributeSafe"
	case FullyTrusted:
		return "FullyTrusted"
	default:
		panic("unknown trust value")
	}
}

type SafeString struct {
	value string
	trust StringTrust
}

func (s SafeString) Convert(target StringTrust) (string, error) {
	if s.trust == FullyTrusted || target == s.trust {
		return s.value, nil
	}

	// Only untrusted values can be sanitized. No level of sanitization can make
	// a value fully trusted for all contexts.
	if s.trust != Untrusted || target == FullyTrusted {
		return "", fmt.Errorf("%w: cannot convert from %v to %v", ErrStringUntrusted, s, target)
	}

	return escape(s.value, target)
}

func (s SafeString) String() string {
	return fmt.Sprintf("%q (%v)", s.value, s.trust)
}

func (s SafeString) Value() string {
	return s.value
}

func TrustedString(s string, trust StringTrust) SafeString {
	return SafeString{s, trust}
}

func FullyTrustedString(s string) SafeString {
	return SafeString{s, FullyTrusted}
}

func UntrustedString(s string) SafeString {
	return SafeString{s, Untrusted}
}

func escape(s string, trust StringTrust) (string, error) {
	switch trust {
	case HTMLSafe:
		return escapeHTML(s)
	case TextSafe:
		return escapeText(s)
	case URLSafe:
		return escapeURL(s)
	case AttributeSafe:
		return escapeAttribute(s)
	default:
		panic("unknown trust level")
	}
}

func escapeHTML(s string) (string, error) {
	return "", errors.New("HTML fragments must currently be FullyTrusted")
}

func escapeText(s string) (string, error) {
	return html.EscapeString(s), nil
}

func escapeURL(s string) (string, error) {
	return safehtml.URLSanitized(s).String(), nil
}

func escapeAttribute(s string) (string, error) {
	// TODO - I am sure this is somehow deficient. Possibly, there is no
	// universal standard to which attributes should be escaped ,and we should
	// always require a fully trusted string?
	return "", errors.New("attributes must currently be FullyTrusted")
}
