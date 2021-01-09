package bindings

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

var ErrUndefined = errors.New("undefined")

type Value struct {
	idx           int
	value         string
	stream        ValueStream
	setError      error
	debugOnlyName string
}

type ValueMap struct {
	Vars    *Map
	values  []string
	streams []ValueStream
}

func (vm *ValueMap) String() string {
	var sb strings.Builder
	vm.DebugDump(&sb, 0)
	return sb.String()
}

func (vm *ValueMap) DebugDump(w io.Writer, depth int) {
	indent := strings.Repeat("\t", depth)

	fmt.Fprintf(w, "%sValueMap{ Strict=%v", indent, vm.Vars.Strict)
	if vm.Vars.nameInParent != "" {
		fmt.Fprintf(w, " (nested, named %q)", vm.Vars.nameInParent)
	}
	fmt.Fprint(w, "\n")

	for i, v := range vm.Vars.vars {
		fmt.Fprintf(w, "%s\tvar %d/%d: %q@%d (%v)\n", indent, i+1, len(vm.Vars.vars), v.name, v.idx, v.level)
		if s := vm.GetString(v); s != "" {
			fmt.Fprintf(w, "%s\t\tstring %q\n", indent, s)
		} else {
			fmt.Fprintf(w, "%s\t\t(empty)\n", indent)
		}
	}

	for i, nm := range vm.Vars.maps {
		fmt.Fprintf(w, "%s\tnested map %d/%d:\n", indent, i+1, len(vm.Vars.maps))
		stream := vm.GetStream(nm)
		if stream == nil {
			fmt.Fprintf(w, "%s\t\t(empty)\n", indent)
			continue
		}

		next := stream.Stream()
		j := 0
		for nvm := next(); nvm != nil; nvm = next() {
			j++
			fmt.Fprintf(w, "%s\tinstantiation (ValueMap) #%d:\n", indent, j)
			nvm.DebugDump(w, depth+2)
		}
	}

	fmt.Fprintf(w, "%s}\n", indent)
}

func (vm *ValueMap) setNestedMapStream(v Value) error {
	limit := len(vm.Vars.maps)
	if limit <= v.idx {
		return fmt.Errorf("%w subsection value stream %s", ErrUndefined, v.debugOnlyName)
	}
	if len(vm.streams) < limit {
		tmp := vm.streams
		vm.streams = make([]ValueStream, limit)
		copy(vm.streams, tmp)
	}

	vm.streams[v.idx] = v.stream
	return nil
}

func (vm *ValueMap) setValue(v Value) error {
	limit := len(vm.Vars.vars)
	if limit <= v.idx {
		return fmt.Errorf("%w var %s", ErrUndefined, v.debugOnlyName)
	}
	if len(vm.values) < limit {
		tmp := vm.values
		vm.values = make([]string, limit)
		copy(vm.values, tmp)
	}

	vm.values[v.idx] = v.value
	return nil
}

func (vm *ValueMap) Set(v Value) error {
	if v.setError != nil {
		return fmt.Errorf("value %q could not be set: %w", v.debugOnlyName, v.setError)
	}
	if v.stream != nil {
		return vm.setNestedMapStream(v)
	}
	return vm.setValue(v)
}

func (vm *ValueMap) GetString(v Var) string {
	if len(vm.values) <= v.idx {
		return ""
	}

	if v.checkOnlyAttachedMap == nil {
		panic(fmt.Sprintf("%v is unattached (programmer error - free variables MUST be attached to a map)", v))
	}
	if v.checkOnlyAttachedMap != vm.Vars {
		panic(fmt.Sprintf("%v is bound to the map %q, this ValueMap is instantiated from %q (programmer error - variable used in wrong context)", v, v.checkOnlyAttachedMap.DebugName(), vm.Vars.DebugName()))
	}

	return vm.values[v.idx]
}

func (vm *ValueMap) GetStream(m *Map) ValueStream {
	if len(vm.streams) <= vm.Vars.idxInParent {
		return nil
	}
	return vm.streams[m.idxInParent]
}

// ValueStream is a collection of ValueMaps used to generate repeated
// subsections of an HTML page. A convenient implementation is ValueSeries, but
// this interface can also be implemented such that values are loaded on the
// fly, to avoid using too much memory when generating large HTML pages.
type ValueStream interface {
	// Stream returns an iterator at position zero. Each call to Stream must
	// return an iterator that yields the same values in the same order.
	Stream() ValueIterator
}

// ValueIterator returns a new ValueSet on each call. It returns nil when
// there are no more values.
type ValueIterator func() *ValueMap

// ValueSeries is a convenient implementation of ValueSetStream, which is
// backed by a slice of ValueSets.
type ValueSeries []*ValueMap

// Stream returns an iterators that returns ValueSets from the slice one at a
// time.
func (s ValueSeries) Stream() ValueIterator {
	f := func() *ValueMap {
		if len(s) == 0 {
			return nil
		}
		vs := s[0]
		s = s[1:]
		return vs
	}
	return ValueIterator(f)
}
