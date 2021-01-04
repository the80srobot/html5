package safe

import (
	"errors"
	"fmt"
	"reflect"
)

// TrustLevel specifies the level of trust in a string value, by its origin or
// construction.
type TrustLevel int16

// ErrStringUntrusted signifies that a SafeString's StringTrust level was
// insufficient in some context. (Check with errors.Is, not with ==.)
var ErrStringUntrusted = errors.New("cannot provide trust level")

// Different StringTrust levels. The order is not significant, except Untrusted,
// which should always be the zero value.
const (
	// Must not appear in the page.
	Untrusted TrustLevel = iota
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

	Default = Untrusted
)

func (l TrustLevel) String() string {
	switch l {
	case Untrusted:
		return "Untrusted (Default)"
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

func Check(s String, required TrustLevel) (string, error) {
	if s.Check(required) {
		return s.String(), nil
	}
	return "", fmt.Errorf("%v %w %v", reflect.TypeOf(s), ErrStringUntrusted, required)
}

func Max(x, y TrustLevel) TrustLevel {
	if x == y {
		return x
	}
	if x == Untrusted {
		return y
	}
	if y == Untrusted {
		return x
	}

	return FullyTrusted
}
