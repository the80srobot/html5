package bindings

import (
	"errors"
	"fmt"
	"io"
	"strings"

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
// BindArg must specify either a string or a ValueStream.
//
// To specify a string, set Value to a non-nil safe.String and Var to a variable
// on which it should be set. If Var is no attached, it will become attached to
// the ValueMap passed to Bind.
//
// To specify a ValueStream, set NestedRows such that each slice of BindArgs
// (each row) contains valid arguments for Bind on the NestedMap. Set NestedMap
// to the corresponding Map.
//
// As a convenience, if Var or NestedMap are not available, but Name is a
// non-empty string, then Bind will lookup or declare the Var or Map
// automatically. This is slower than providing the Var or NestedMap.
//
// Optionally, if TrustRequirement is specified, it will be treated as an extra
// requirement on top of the Var's requirement, and Var will be promoted to
// safe.Max of the two.
type BindArg struct {
	// Fields to specify a safe string value.
	Var   Var
	Value safe.String

	// Fields to specify a nested ValueStream.
	NestedMap  *Map
	NestedRows [][]BindArg

	// Optional: if Var/NestedMap are not specified, lookup or declare using
	// this name.
	Name string
	// Optional: if specified, Var will be promoted to safe.Max of this and its
	// existing trust requirement.
	TrustRequirement safe.TrustLevel
}

func (a BindArg) String() string {
	var sb strings.Builder

	sb.WriteString("BindArg{ ")
	a.describe(&sb)
	sb.WriteString("}")
	return sb.String()
}

func (a BindArg) describe(w io.Writer) {
	if a.Name != "" {
		fmt.Fprintf(w, "named %q ", a.Name)
	}

	if a.Value != nil {
		fmt.Fprintf(w, "var=%v, string=%v, req=%v ", a.Var, a.Value, a.TrustRequirement)
	}

	if a.NestedMap != nil {
		fmt.Fprintf(w, "map named %q, idx=%d, ", a.NestedMap.DebugName(), a.NestedMap.idxInParent)
	}

	if a.NestedRows != nil {
		fmt.Fprintf(w, "ValueStream of %d rows", len(a.NestedRows))
	}
}

// DebugDump writes a verbose description of the BindArg to the writer. Depth
// can be used to indent output with tabs.
func (a BindArg) DebugDump(w io.Writer, depth int) {
	indent := strings.Repeat("  ", depth)
	io.WriteString(w, indent)
	io.WriteString(w, "BindArg{")
	a.describe(w)

	for i, row := range a.NestedRows {
		fmt.Fprintf(w, "\n%s  row %d/%d:\n", indent, i+1, len(a.NestedRows))
		for _, col := range row {
			col.DebugDump(w, depth+2)
			io.WriteString(w, "\n")
		}
	}
	if len(a.NestedRows) != 0 {
		fmt.Fprintf(w, "\n%s", indent)
	}
	io.WriteString(w, "}")
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

// BindArgs is a slice of BindArg values. Its only advantage over a plain slice
// is that it knows how to DebugDump itself, and so can be passed to functions
// that accept DebugDumper.
type BindArgs []BindArg

// DebugDump writes a verbose description of the BindArgs to the writer.
func (a BindArgs) DebugDump(w io.Writer, depth int) {
	indent := strings.Repeat("  ", depth)
	for i, v := range a {
		fmt.Fprintf(w, "%sarg %d/%d: ", indent, i+1, len(a))
		v.DebugDump(w, depth)
		io.WriteString(w, "\n")
	}
}
