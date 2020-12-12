package html

import (
	"errors"
	"testing"
)

func TestSafeString(t *testing.T) {
	for _, tc := range []struct {
		comment      string
		input        SafeString
		convertTrust StringTrust
		wantString   string
		wantErr      error
	}{
		{
			comment:      "trusted - trusted",
			convertTrust: FullyTrusted,
			input:        FullyTrustedString("<p>Foo!</p>"),
			wantString:   "<p>Foo!</p>",
			wantErr:      nil,
		},
		{
			comment:      "untrusted - trusted",
			convertTrust: FullyTrusted,
			input:        UntrustedString("<p>Foo!</p>"),
			wantString:   "",
			wantErr:      ErrStringUntrusted,
		},
		{
			comment:      "untrusted - text",
			convertTrust: TextSafe,
			input:        UntrustedString("<p>Foo!</p>"),
			wantString:   `&lt;p&gt;Foo!&lt;/p&gt;`,
			wantErr:      nil,
		},
	} {
		t.Run(tc.comment, func(t *testing.T) {
			out, err := tc.input.Convert(tc.convertTrust)
			if !errors.Is(err, tc.wantErr) || out != tc.wantString {
				t.Errorf("%v.Convert(%v) => %q, %v (wanted %q, %v)",
					tc.input, tc.convertTrust, out, err, tc.wantString, tc.wantErr)
			}
		})
	}
}
