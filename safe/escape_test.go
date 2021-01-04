package safe

import (
	"errors"
	"testing"
)

func TestEscapeURL(t *testing.T) {
	for _, tc := range []struct {
		comment string
		input   string
		want    string
		wantErr error
	}{
		{
			comment: "valid http",
			input:   "http://google.com/",
			want:    "http://google.com/",
		},
		{
			comment: "injection",
			input:   `http://" href="evil.com`,
			wantErr: errInvalidInput,
		},
		{
			comment: "fragment",
			input:   "/articles",
			want:    "/articles",
		},
		{
			comment: "bad schema",
			input:   "ssh://adam@foo.local",
			wantErr: errForbiddenSchema,
		},
	} {
		t.Run(tc.comment, func(t *testing.T) {
			// This almost tests the exported API, but the unexported function
			// makes it easier to check for specific errors.
			//
			// TODO(adam): Should we expose the errors?
			s, err := escapeURL(tc.input)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("EscapeURL(%v) => (%v, %v), wanted error %v", tc.input, s, err, tc.wantErr)
			}
			if s != tc.want {
				t.Errorf("EscapeURL(%v) => (%v, %v), wanted %v", tc.input, s, err, tc.want)
			}
		})
	}
}
