package safe

import (
	"errors"
)

type safestring string

type URL struct {
	s string
}

type Attribute struct {
	s string
}

type HTML struct {
	s string
}

type Text struct {
	s string
}

// Level specifies the level of trust in a string value, by its origin or
// construction.
type Level int16

// ErrStringUntrusted signifies that a SafeString's StringTrust level was
// insufficient in some context. (Check with errors.Is, not with ==.)
var ErrStringUntrusted = errors.New("cannnot satisfy trust level")

// Different StringTrust levels. The order is not significant, except Untrusted,
// which should always be the zero value.
const (
	// Must not appear in the page.
	Untrusted Level = iota
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

func (l Level) String() string {
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

type String interface {
	Safe(Level) string
}

type constant string

func Bless(s string) constant {
	return constant(s)
}

func (c constant) Safe(Level) string {
	return string(c)
}

type Raw string
