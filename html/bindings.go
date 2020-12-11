package html

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// Tag identifies a binding (string or subsection). BindingSet and ValueSet use
// Tags to lookup bindings and values without having to hash their names.
//
// A Tag can refer to a string binding, or a subsection binding, but not both.
// Attempting to use a string tag to lookup a subsection will result in an
// error.
type Tag struct {
	v int
}

// ZeroTag is the default value of Tag. It is neither a valid string tag, nor a
// valid subsection tag.
var ZeroTag = Tag{}

func stringTag(i int) Tag {
	return Tag{i + 1}
}

func subsectionTag(i int) Tag {
	return Tag{-i - 1}
}

func (t Tag) string() int {
	return t.v - 1
}

func (t Tag) subsection() int {
	return t.v + 1
}

func (t Tag) findString(v []StringBinding) (int, error) {
	if t.v < 1 {
		return -1, fmt.Errorf("invalid string tag %v", t)
	}

	i := t.v - 1
	if len(v) <= i {
		return -1, fmt.Errorf("string tag %v not found", t)
	}
	return i, nil
}

func (t Tag) findSubsection(v []*BindingSet) (int, error) {
	if t.v > -1 {
		return -1, fmt.Errorf("invalid subsection tag %v", t)
	}

	i := -t.v - 1
	if len(v) <= i {
		return -1, fmt.Errorf("subsection tag %v not found", t)
	}
	return i, nil
}

func (t Tag) String() string {
	switch {
	case t == ZeroTag:
		return "Tag{0}"
	case t.v < 0:
		return fmt.Sprintf("Tag{section %d}", -t.v)
	default:
		return fmt.Sprintf("Tag{string %d}", t.v)
	}
}

// BindingSet is a collection of the dynamic properties of a page (bindings).
// Each binding is either a string, or another BindingSet describing the dynamic
// properties of a subsection.
//
// The BindingSet contains only definitions, and is usually created together
// with a page Template. Call Bind to create a ValueSet for the bindings in this
// BindingSet.
type BindingSet struct {
	strings     []StringBinding
	stringNames map[string]Tag

	// Binding sets for subsections, which can be repeated in the ValueSet.
	subsections     []*BindingSet
	subsectionNames map[string]Tag
}

func (bs *BindingSet) lazyInit() {
	if len(bs.strings) == 0 && len(bs.subsections) == 0 {
		bs.stringNames = make(map[string]Tag)
		bs.subsectionNames = make(map[string]Tag)
	}
}

// AddString creates a string binding at the given level of trust. It returns a
// Tag, which can later be used to provide a value for this binding.
func (bs *BindingSet) AddString(name string, trust StringTrust) Tag {
	bs.lazyInit()

	tag, ok := bs.stringNames[name]
	if !ok {
		idx := len(bs.strings)
		tag = stringTag(idx)
		bs.strings = append(bs.strings, StringBinding{Name: name, Trust: trust, tag: tag})
		bs.stringNames[name] = tag
	}
	return tag
}

// AddSubsection creates a subsection binding and returns a tag used to later
// assign values to it.
func (bs *BindingSet) AddSubsection(name string, subsectionBindings *BindingSet) Tag {
	bs.lazyInit()

	tag, ok := bs.subsectionNames[name]
	if !ok {
		idx := len(bs.subsections)
		tag = subsectionTag(idx)
		bs.subsections = append(bs.subsections, subsectionBindings)
		bs.subsectionNames[name] = tag
	}
	return tag
}

// StringTag looks up a string binding by its name, and returns the tag used to
// identify it. It's better to save the tag returned from AddString, which is
// the same value. If the set contains no binding with the given name, then the
// tag will be TagZero.
func (bs *BindingSet) StringTag(name string) Tag {
	return bs.stringNames[name]
}

// SubsectionTag is the same as StringTag, but works on subsection bindings.
func (bs *BindingSet) SubsectionTag(name string) Tag {
	return bs.subsectionNames[name]
}

// Bind creates a ValueSet from the provided binding values. Each Value must
// refer to a binding in this BindingSet using a valid tag or name, otherwise
// Bind will return an error. Prefer to identify bindings by their tags, to
// avoid multiple hash table lookups.
func (bs *BindingSet) Bind(values ...ValueArg) (*ValueSet, error) {
	vs := ValueSet{BindingSet: bs}
	for _, v := range values {
		if err := vs.Bind(&v); err != nil {
			return nil, err
		}
	}
	return &vs, nil
}

// ValueArg specifies value of a binding. It is a convenient way of supplying
// arguments to BindingSet.Bind, but has no other use.
type ValueArg struct {
	// Identity of the binding. A string binding must be used with StringValue
	// and a subsection binding with a subsection value, and the two are
	// mutually exclusive.
	//
	// Prefer to use Tag over Name, as it is faster.
	Tag  Tag
	Name string

	StringValue string
	StringTrust StringTrust

	Subsections [][]ValueArg
}

func (v ValueArg) String() string {
	return fmt.Sprintf("ValueArg{tag=%v, name=%s, trust=%v, value=%q, subsection=%v}", v.Tag, v.Name, v.StringTrust, v.StringValue, v.Subsections)
}

// ValueSetStream is a collection of ValueSets used to generate repeated
// subsections of an HTML page. A convenient implementation is ValueSetSlice,
// but this interface can also be implemented such that values are loaded on the
// fly, to avoid using too much memory when generating large HTML pages.
type ValueSetStream interface {
	// Stream returns an iterator at position zero. Each call to Stream must
	// return an iterator that yields the same values in the same order.
	Stream() ValueSetIterator
}

// ValueSetIterator returns a new ValueSet on each call. It returns nil when
// there are no more values.
type ValueSetIterator func() *ValueSet

// ValueSetSlice is a convenient implementation of ValueSetStream, which is
// backed by a slice of ValueSets.
type ValueSetSlice []*ValueSet

// Stream returns an iterators that returns ValueSets from the slice one at a
// time.
func (s ValueSetSlice) Stream() ValueSetIterator {
	f := func() *ValueSet {
		if len(s) == 0 {
			return nil
		}
		vs := s[0]
		s = s[1:]
		return vs
	}
	return ValueSetIterator(f)
}

// ValueSet is a collection of dynamic values used to generate an HTML page.
// Each ValueSet may contain only values defined by its BindingSet, but it is
// not required to have a value for every binding in the BindingSet.
type ValueSet struct {
	BindingSet  *BindingSet
	strings     []string
	subsections []ValueSetStream
}

func (vs *ValueSet) String() string {
	var sb strings.Builder
	sb.WriteString("ValueSet{\n")
	for name, tag := range vs.BindingSet.stringNames {
		fmt.Fprintf(&sb, "  binding %s (tag %v): ", name, tag)
		if err := vs.writeStringTo(&sb, tag); err != nil {
			sb.WriteString("<nil>")
		}
		sb.WriteByte('\n')
	}
	sb.WriteByte('}')
	return sb.String()
}

// Bind will set the value of a single binding, provided its definition is found
// in the BindingSet. See ValueArg for more.
func (vs *ValueSet) Bind(v *ValueArg) error {
	switch {
	case v.StringValue != "":
		tag := v.Tag
		if tag == ZeroTag {
			tag = vs.BindingSet.StringTag(v.Name)
		}
		return vs.BindString(tag, v.StringTrust, v.StringValue)
	case len(v.Subsections) != 0:
		tag := v.Tag
		if tag == ZeroTag {
			tag = vs.BindingSet.SubsectionTag(v.Name)
		}
		return vs.BindSubsections(tag, v.Subsections)
	default:
		return errors.New("empty value")
	}
}

// BindString will set the value of a single string binding, provided its
// defintion is found in the BindingSet. See StringTrust for a discussion of
// string safety.
func (vs *ValueSet) BindString(tag Tag, valueTrust StringTrust, value string) error {
	idx, err := tag.findString(vs.BindingSet.strings)
	if err != nil {
		return err
	}

	b := vs.BindingSet.strings[idx]
	s, err := b.SafeString(valueTrust, value)
	if err != nil {
		return fmt.Errorf("SafeString for %v: %w", tag, err)
	}

	if len(vs.strings) < len(vs.BindingSet.strings) {
		values := vs.strings
		vs.strings = make([]string, len(vs.BindingSet.strings))
		copy(vs.strings, values)
	}
	vs.strings[idx] = s
	return nil
}

// BindSubsections will create a collection of nested ValueSets for the given
// subsection. If there is no such subsection in the BindingSet, this will
// return an error. Each value slice must contain valid arguments for
// BindingSet.Bind on the subsection's BindingSet.
func (vs *ValueSet) BindSubsections(tag Tag, subsectionValues [][]ValueArg) error {
	idx, err := tag.findSubsection(vs.BindingSet.subsections)
	if err != nil {
		return err
	}
	subsectionBindings := vs.BindingSet.subsections[idx]

	var subsections ValueSetSlice
	for _, values := range subsectionValues {
		subsection, err := subsectionBindings.Bind(values...)
		if err != nil {
			return err
		}
		subsections = append(subsections, subsection)
	}
	return vs.BindSubsectionStream(tag, subsections)
}

// BindSubsectionStream will set the value stream for the specified subsection.
// If there is no such subsection in the BindingSet, this will return an error.
// The value stream must yield values valid for BindingSet.Bind on the
// subsection's BindingSet.
func (vs *ValueSet) BindSubsectionStream(tag Tag, stream ValueSetStream) error {
	idx, err := tag.findSubsection(vs.BindingSet.subsections)
	if err != nil {
		return err
	}

	if len(vs.subsections) < len(vs.BindingSet.subsections) {
		values := vs.subsections
		vs.subsections = make([]ValueSetStream, len(vs.BindingSet.subsections))
		copy(vs.subsections, values)
	}
	vs.subsections[idx] = stream
	return nil
}

func (vs *ValueSet) writeStringTo(w io.Writer, tag Tag) error {
	idx, err := tag.findString(vs.BindingSet.strings)
	if err != nil {
		return err
	}
	if len(vs.strings) <= idx {
		// Valid tag, but no value provided. Same as an empty string.
		return nil
	}
	_, err = io.WriteString(w, vs.strings[idx])
	return err
}

func (vs *ValueSet) iterateSubsection(tag Tag) (ValueSetIterator, error) {
	idx, err := tag.findSubsection(vs.BindingSet.subsections)
	if err != nil {
		return nil, err
	}
	if len(vs.subsections) <= idx {
		// Valid tag, but no subsections. Return an empty slice.
		return ValueSetSlice(nil).Stream(), nil
	}
	return vs.subsections[idx].Stream(), nil
}
