package bindings

import (
	"errors"

	"github.com/the80srobot/html5/safe"
)

// GetStringByName looks up the value of a Var with the given name. Doesn't
// distinguish between empty string and no value. If the Var was previously
// undeclared, then this will have the side effect of declaring it in the
// associated Map.
//
// This is slower than ValueMap.GetString and should only be used for debugging.
func GetStringByName(vm *ValueMap, name string) string {
	v := vm.Vars.Declare(name, safe.Default)
	return vm.GetString(v)
}

// Bind applies the provided bindings to the ValueMap. This is the same as
// calling ValueMap.Set repeatedly, but allows the caller to conveniently
// specify entire nested structures in one call.
//
// See BindArg for more.
func Bind(vm *ValueMap, args ...BindArg) error {
	for _, arg := range args {
		if arg.Value == nil {
			if err := bindSubsection(vm, arg); err != nil {
				return err
			}
			continue
		}

		if err := bindString(vm, arg); err != nil {
			return err
		}
	}
	return nil
}

// BindArg specifies a single string or nested map to set on a ValueMap. Each
// BindArg must specify either a string or a ValueStream (as more BindArgs).
//
// If a BindArg specifies a string, then Var and Value must be set.
// TrustRequirement may be left at Default, or set to a new requirement, in
// which case the underlying Map will increase the required trust for this var
// to safe.Max of the previous trust and the new value.
//
// If a BindArg specifies a nested ValueStream, then NestedMap must be set to
// the Map, and NestedRows must be set to a slice of slices of BindArgs, such
// that a nested ValueMap can be built recursively from each row.
type BindArg struct {
	Name             string
	TrustRequirement safe.TrustLevel
	Var              Var
	NestedMap        *Map

	Value      safe.String
	NestedRows [][]BindArg
}

func bindSubsection(vm *ValueMap, arg BindArg) error {
	m := arg.NestedMap
	if m == nil {
		m = vm.Vars.Nest(arg.Name)
	} else {
		m2 := vm.Vars.Nest(m.nameInParent)
		if m2 != m {
			return errors.New("cannot use a map that's not attached to this parent")
		}
	}

	var series ValueSeries
	for _, row := range arg.NestedRows {
		nestedVM := m.MustBind()
		for _, col := range row {
			if err := Bind(nestedVM, col); err != nil {
				return err
			}
		}
		series = append(series, nestedVM)
	}
	return vm.Set(m.BindStream(series))
}

func bindString(vm *ValueMap, arg BindArg) error {
	v := arg.Var
	if v == ZeroVar {
		v = vm.Vars.Declare(arg.Name, arg.TrustRequirement)
	} else {
		v = vm.Vars.Attach(v, arg.TrustRequirement)
	}

	return vm.Set(v.Bind(arg.Value))
}
