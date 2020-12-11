package html

import (
	"errors"
	"testing"
)

func TestSafeString(t *testing.T) {
	for _, tc := range []struct {
		comment    string
		binding    StringBinding
		input      string
		inputTrust StringTrust
		wantString string
		wantErr    error
	}{
		{
			comment:    "trusted - trusted",
			binding:    StringBinding{Trust: FullyTrusted},
			input:      "<p>Foo!</p>",
			inputTrust: FullyTrusted,
			wantString: "<p>Foo!</p>",
			wantErr:    nil,
		},
		{
			comment:    "untrusted - trusted",
			binding:    StringBinding{Trust: FullyTrusted},
			input:      "<p>Foo!</p>",
			inputTrust: Untrusted,
			wantString: "",
			wantErr:    ErrStringUntrusted,
		},
		{
			comment:    "untrusted - text",
			binding:    StringBinding{Trust: TextSafe},
			input:      "<p>Foo!</p>",
			inputTrust: Untrusted,
			wantString: `&lt;p&gt;Foo!&lt;/p&gt;`,
			wantErr:    nil,
		},
	} {
		t.Run(tc.comment, func(t *testing.T) {
			out, err := tc.binding.SafeString(tc.inputTrust, tc.input)
			if !errors.Is(err, tc.wantErr) || out != tc.wantString {
				t.Errorf("%v.SafeString(%v, %q) => %q, %v (wanted %q, %v)",
					tc.binding, tc.inputTrust, tc.input, out, err, tc.wantString, tc.wantErr)
			}
		})
	}
}
