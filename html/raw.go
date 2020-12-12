package html

type RawNode struct {
	HTML SafeString
}

func (r *RawNode) compile(tc *templateCompiler, _ int, _ *CompileOptions) error {
	s, err := r.HTML.Convert(FullyTrusted)
	if err != nil {
		return err
	}
	_, err = tc.WriteString(s)
	return err
}
