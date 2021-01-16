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
	return fmt.Sprintf("Var{%d, %q, %v}", v.idx, v.name, v.level)
}

func (v Var) Check(required safe.TrustLevel) bool {
	return true
}

var ZeroVar = Var{}

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

func (v Var) Bind(ss safe.String) Value {
	value, err := v.tryBind(ss)
	value.setError = err
	return value
}

func (v Var) Attached() bool {
	return v.checkOnlyAttachedMap != nil
}

type constString string

func (v Var) BindConst(value constString) Value {
	return Value{debugOnlyName: v.name, idx: v.idx, value: string(value), checkOnlyContainingMap: v.checkOnlyAttachedMap}
}

type Map struct {
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

func (m *Map) Root() bool {
	return m.checkOnlyParentMap == nil
}

func (m *Map) String() string {
	var sb strings.Builder
	m.DebugDump(&sb, 0)
	return sb.String()
}

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

func (m *Map) Attach(v Var, level safe.TrustLevel) Var {
	return m.Declare(v.name, safe.Max(v.level, level))
}

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

func (m *Map) MustBind(values ...Value) *ValueMap {
	vm, err := m.Bind(values...)
	if err != nil {
		panic(err)
	}
	return vm
}

func (m *Map) BindStream(stream ValueStream) Value {
	return Value{
		idx:                    m.idxInParent,
		stream:                 stream,
		debugOnlyName:          m.nameInParent,
		checkOnlyContainingMap: m.checkOnlyParentMap,
	}
}

func (m *Map) BindSeries(maps ...*ValueMap) Value {
	return m.BindStream(ValueSeries(maps))
}

func (m *Map) DebugName() string {
	if m.Root() {
		return "root"
	}
	return m.nameInParent
}

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
