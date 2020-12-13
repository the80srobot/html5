package html

// MultiNode concatenates several other nodes.
type MultiNode struct {
	Contents []Node
}

func (m *MultiNode) compile(db *templateCompiler, depth int, opts *CompileOptions) error {
	for _, c := range m.Contents {
		if err := c.compile(db, depth, opts); err != nil {
			return err
		}
	}
	return nil
}
