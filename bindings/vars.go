package bindings

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/the80srobot/html5/safe"
)

// Var uniquelly names a string variable and its required trust level.
type Var struct {
	idx   int
	name  string
	level safe.TrustLevel
	// The map this var is in. Only to be used for correctness checks.
	checkOnlyAttachedMap *Map
}

func (v Var) String() string {
	if v == ZeroVar {
		return "ZeroVar"
	}
	return fmt.Sprintf("Var{%d, %q, %v}", v.idx, v.name, v.level)
}

// Check whether this Var's trust level can satisfy the required trust level.
func (v Var) Check(required safe.TrustLevel) bool {
	return v.level == required ||
		v.level == safe.FullyTrusted ||
		required == safe.Untrusted ||
		(v.level == safe.HTMLSafe && required == safe.TextSafe)
}

// ZeroVar is a an empty Var.
var ZeroVar = Var{}

// Declare an unnatached Var as a placeholder. Unnatached Vars cannot be used
// with ValueMaps, but can be converted to attached Vars using Map.Attach.
func Declare(name string) Var {
	if name == "" {
		panic("Var name cannot be empty")
	}
	return Var{name: name, level: safe.Untrusted}
}

func (v Var) tryBind(ss safe.String) (Value, error) {
	s, err := safe.Check(ss, v.level)
	if err != nil {
		return Value{}, fmt.Errorf("binding value %s: %w", v.name, err)
	}
	return Value{debugOnlyName: v.name, idx: v.idx, value: s, checkOnlyContainingMap: v.checkOnlyAttachedMap}, nil
}

// Bind returns a Value created by binding the provided string to this Var.
func (v Var) Bind(ss safe.String) Value {
	value, err := v.tryBind(ss)
	if err != nil {
		return Value{
			debugOnlyName: v.name,
			trustErr:      err,
		}
	}
	return value
}

// Attached returns whether this Var is attached to a Map.
func (v Var) Attached() bool {
	return v.checkOnlyAttachedMap != nil
}

type constString string

// BindConst is like Bind, but bypasses the trust level check. To guarantee
// safety, BindConst only accepts string literals.
func (v Var) BindConst(value constString) Value {
	return Value{debugOnlyName: v.name, idx: v.idx, value: string(value), checkOnlyContainingMap: v.checkOnlyAttachedMap}
}

// Map is a collection of Vars and nested Maps. Each html5.Template uses a
// single root Map to hold all the dynamic elements of the page, and their
// requisite levels of trust.
//
// Maps are "instantiated" into ValueMaps, which specify a set of values for the
// Vars and nested Maps in the Map.
type Map struct {
	// If true, then the Map won't accept Bind arguments for non-existent Vars.
	Strict bool

	idxInParent  int
	nameInParent string
	vars         []Var
	varsByName   map[string]int
	maps         []*Map
	mapsByName   map[string]int

	disallowCopy       sync.Mutex
	checkOnlyParentMap *Map
}

// Root returns whether this Map is the top-most Map in the hierarchy, or a
// nested Map.
func (m *Map) Root() bool {
	return m.checkOnlyParentMap == nil
}

func (m *Map) String() string {
	var sb strings.Builder
	m.DebugDump(&sb, 0)
	return sb.String()
}

// Declare a Var with the given name, at the given trust level.
func (m *Map) Declare(name string, level safe.TrustLevel) Var {
	if name == "" {
		panic("Var name cannot be empty")
	}

	idx, ok := m.varsByName[name]
	if ok {
		m.vars[idx].level = safe.Max(m.vars[idx].level, level)
		return m.vars[idx]
	}

	idx = len(m.vars)
	m.vars = append(m.vars, Var{idx: idx, name: name, level: level, checkOnlyAttachedMap: m})
	if m.varsByName == nil {
		m.varsByName = map[string]int{name: idx}
	} else {
		m.varsByName[name] = idx
	}
	return m.vars[idx]
}

// Attach returns a copy of the provided free Var that's associated to this Map.
func (m *Map) Attach(v Var, level safe.TrustLevel) Var {
	return m.Declare(v.name, safe.Max(v.level, level))
}

// Nest creates a nested Map with the given name and returns it. Nested Maps can
// be used to create ValueStreams, which specify repeated sections in the HTML
// page. (For example, comments under an article.) Maps can be nested to
// arbitrary depth.
func (m *Map) Nest(name string) *Map {
	idx, ok := m.mapsByName[name]
	if ok {
		return m.maps[idx]
	}

	idx = len(m.maps)
	m.maps = append(m.maps, &Map{idxInParent: idx, nameInParent: name, checkOnlyParentMap: m})
	if m.mapsByName == nil {
		m.mapsByName = map[string]int{name: idx}
	} else {
		m.mapsByName[name] = idx
	}
	return m.maps[idx]
}

// Bind creates a ValueMap and sets the provided values, if any. The Values must
// be associated to Vars associated to this Map, otherwise Bind will panic.
func (m *Map) Bind(values ...Value) (*ValueMap, error) {
	vm := &ValueMap{
		Vars: m,
	}

	for _, value := range values {
		err := vm.Set(value)
		if !m.Strict && errors.Is(err, ErrUndefined) {
			err = nil
		}
		if err != nil {
			return nil, err
		}
	}

	return vm, nil
}

// MustBind is like Bind, but panics on error.
func (m *Map) MustBind(values ...Value) *ValueMap {
	vm, err := m.Bind(values...)
	if err != nil {
		panic(err)
	}
	return vm
}

// BindStream binds the provided ValueStream to this Map, returning a Value that
// can be set on a parent ValueMap.
func (m *Map) BindStream(stream ValueStream) Value {
	return Value{
		idx:                    m.idxInParent,
		stream:                 stream,
		debugOnlyName:          m.nameInParent,
		checkOnlyContainingMap: m.checkOnlyParentMap,
	}
}

// BindSeries is like BindStream, but constructs a ValueStream for the caller.
func (m *Map) BindSeries(maps ...*ValueMap) Value {
	return m.BindStream(ValueSeries(maps))
}

// DebugName returns a name of this map, or "root" if it isn't nested. Used for
// debugging only.
func (m *Map) DebugName() string {
	if m == nil {
		return "<nil>"
	}

	if m.Root() {
		return "root"
	}
	return m.nameInParent
}

// DebugDump writes a detailed description of this Map to the writer.
func (m *Map) DebugDump(w io.Writer, depth int) {
	indent := strings.Repeat("\t", depth)

	fmt.Fprintf(w, "%sMap{ Strict=%v ", indent, m.Strict)
	if m.nameInParent != "" {
		fmt.Fprintf(w, "(nested, named %q)", m.nameInParent)
	}
	fmt.Fprint(w, "\n")

	for i, v := range m.vars {
		fmt.Fprintf(w, "%s\tvar %d/%d: %q@%d (%v)\n", indent, i+1, len(m.vars), v.name, v.idx, v.level)
	}

	for i, nm := range m.maps {
		fmt.Fprintf(w, "%s\tnested map %d/%d:\n", indent, i+1, len(m.maps))
		nm.DebugDump(w, depth+1)
	}

	fmt.Fprintf(w, "%s}\n", indent)
}
