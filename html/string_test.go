package html

import (
	"errors"
	"testing"
)

func TestSafeString(t *testing.T) {
	for _, tc := range []struct {
		comment    string
		binding    StringBinding
		input      SafeString
		inputTrust StringTrust
		wantString string
		wantErr    error
	}{
		{
			comment:    "trusted - trusted",
			binding:    StringBinding{Trust: FullyTrusted},
			input:      FullyTrustedString("<p>Foo!</p>"),
			wantString: "<p>Foo!</p>",
			wantErr:    nil,
		},
		{
			comment:    "untrusted - trusted",
			binding:    StringBinding{Trust: FullyTrusted},
			input:      UntrustedString("<p>Foo!</p>"),
			wantString: "",
			wantErr:    ErrStringUntrusted,
		},
		{
			comment:    "untrusted - text",
			binding:    StringBinding{Trust: TextSafe},
			input:      UntrustedString("<p>Foo!</p>"),
			wantString: `&lt;p&gt;Foo!&lt;/p&gt;`,
			wantErr:    nil,
		},
	} {
		t.Run(tc.comment, func(t *testing.T) {
			out, err := tc.binding.Convert(tc.input)
			if !errors.Is(err, tc.wantErr) || out != tc.wantString {
				t.Errorf("%v.SafeString(%v, %q) => %q, %v (wanted %q, %v)",
					tc.binding, tc.inputTrust, tc.input, out, err, tc.wantString, tc.wantErr)
			}
		})
	}
}
