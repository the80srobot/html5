package bindings

import (
	"errors"

	"github.com/the80srobot/html5/safe"
)

func GetStringByName(vm *ValueMap, name string) string {
	v := vm.Vars.Declare(name, safe.Default)
	return vm.GetString(v)
}

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
