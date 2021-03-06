package html5

import (
	"fmt"
	"io"
	"strings"

	"github.com/the80srobot/html5/bindings"
)

type Node interface {
	Content
	compile(tc *templateCompiler, depth int, opts *CompileOptions) error
}

type Template struct {
	Bindings *bindings.Map
	chunks   []chunk
}

func GenerateHTML(w io.Writer, t *Template, values ...bindings.BindArg) error {
	vm, err := t.Bindings.Bind()
	if err != nil {
		return err
	}
	if err := bindings.Bind(vm, values...); err != nil {
		return err
	}

	return t.GenerateHTML(w, vm)
}

func (t *Template) GenerateHTML(w io.Writer, vm *bindings.ValueMap) error {
	for i, chunk := range t.chunks {
		if err := chunk.build(w, vm); err != nil {
			return fmt.Errorf("building chunk #%d of %d: %w", i, len(t.chunks), err)
		}
	}
	return nil
}

func (t *Template) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Template{\n")
	for i, c := range t.chunks {
		fmt.Fprintf(&sb, "\tchunk %d/%d: %v\n", i+1, len(t.chunks), c)
	}
	sb.WriteString("\n\t-- bindings follow after this line --\n\n")
	t.Bindings.DebugDump(&sb, 1)
	sb.WriteByte('}')
	return sb.String()
}

func Compile(n Node, m *bindings.Map, opts *CompileOptions) (*Template, error) {
	tc := &templateCompiler{bindings: m}
	tc.separateChunks = opts.SeparateStaticChunks
	if err := n.compile(tc, opts.RootDepth, opts); err != nil {
		return nil, err
	}
	tc.flush()
	return &Template{chunks: tc.chunks, Bindings: tc.bindings}, nil
}

func MustCompile(n Node, m *bindings.Map, opts *CompileOptions) *Template {
	t, err := Compile(n, m, opts)
	if err != nil {
		panic(err)
	}
	return t
}

type chunk interface {
	build(w io.Writer, vm *bindings.ValueMap) error
}

type staticChunk struct {
	data string
}

func (sc staticChunk) build(w io.Writer, _ *bindings.ValueMap) error {
	_, err := io.WriteString(w, sc.data)
	return err
}

func (sc staticChunk) String() string {
	return fmt.Sprintf("static{%q}", sc.data)
}

type stringBindingChunk struct {
	binding bindings.Var
}

func (sbc stringBindingChunk) build(w io.Writer, vm *bindings.ValueMap) error {
	_, err := io.WriteString(w, vm.GetString(sbc.binding))
	return err

}

func (sbc stringBindingChunk) String() string {
	return fmt.Sprintf("stringBinding{%v}", sbc.binding)
}
