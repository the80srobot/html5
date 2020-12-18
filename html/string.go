package html

import (
	"errors"
	"fmt"
	"html"

	"github.com/google/safehtml"
)

// StringTrust specifies the level of trust in a string value, by its origin or
// construction.
type StringTrust int16

// ErrStringUntrusted signifies that a SafeString's StringTrust level was
// insufficient in some context. (Check with errors.Is, not with ==.)
var ErrStringUntrusted = errors.New("cannnot satisfy trust level")

// Different StringTrust levels. The order is not significant, except Untrusted,
// which should always be the zero value.
const (
	// Must not appear in the page.
	Untrusted StringTrust = iota
	// Safe to insert as HTML.
	HTMLSafe
	// Safe to insert as text (no HTML tags present).
	TextSafe
	// Safe to insert into url attributes (such as href).
	URLSafe
	// Safe to insert into most attributes.
	AttributeSafe
	// Assumed, by origin, safe in any context.
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

// SafeString keeps track of a single string value used to generate a page. It
// can be either constant or a binding. Each string value comes with a trust
// level that specifies where the string is safe for use. The trust level must
// be provided when constructing a constant string, or when assigning a string
// value to a binding.
type SafeString struct {
	binding string
	value   string
	trust   StringTrust
}

// DebugValue returns the string's value as is. Only for logging and debugging.
// Us Convert for proper access.
func (s SafeString) DebugValue() string {
	return s.value
}

// Constant returns whether this string has a constant value. (The other option
// is that the string is a binding, whose value will be supplied later.)
func (s SafeString) Constant() bool {
	return s.binding == ""
}

// Convert tries return the value of this string escaped or trusted for use at
// the target trust level. Fails if no conversion is possible (such as to
// FullyTrusted), or if the string has no constant value.
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
	if s.Constant() {
		return fmt.Sprintf("%q (%v)", s.value, s.trust)
	}
	return fmt.Sprintf("string binding %s (default %q %v)", s.binding, s.value, s.trust)
}

// Binding creates a SafeString binding with the given name. (The value and
// corresponding trust level must be supplied through BindingSet.Binding.)
func Binding(name string) SafeString {
	return SafeString{binding: name}
}

// TrustedString creates a SafeString from the provided value and trust level.
// The value must be known to the caller, by origin or construction, to be safe
// at the specified trust level.
func TrustedString(s string, trust StringTrust) SafeString {
	return SafeString{value: s, trust: trust}
}

// FullyTrustedString creates a fully trusted SafeString from the provided
// value. The value must be known to the caller, by origin or construction, to
// be safe for use in any context. (Usually fully trusted values are hardcoded
// constants that require no escaping or further care.)
func FullyTrustedString(s string) SafeString {
	return SafeString{value: s, trust: FullyTrusted}
}

// UntrustedString creates an untrusted SafeString from the provided value.
// Untrusted strings are not permitted to appear in the page without escaping.
// SafeString.Convert can be used to obtain values at a higher trust level.
func UntrustedString(s string) SafeString {
	return SafeString{value: s, trust: Untrusted}
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
	// TODO - this is the only place we depend on an outside package. It also
	// doesn't really escape the URL, although it does validate it.
	return safehtml.URLSanitized(s).String(), nil
}

func escapeAttribute(s string) (string, error) {
	return "", errors.New("attributes must currently be FullyTrusted")
}
