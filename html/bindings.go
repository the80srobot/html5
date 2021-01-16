package html

import (
	"github.com/the80srobot/html5/safe"
)

// Value is a placeholder for a string value. The two types of values used in
// this package are binding.Var and safe.Value. The difference between the two
// is that Vars can be bound to a specific Value at a later time, while Strings
// are interpolated when we compile the template.
//
// Value provides a way to check that the var or safe string conform to a
// minimal level of trust.
type Value interface {
	Check(required safe.TrustLevel) bool
}
