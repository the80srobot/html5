package safe

import (
	"errors"
	"testing"
)

func TestCheck(t *testing.T) {
	for _, tc := range []struct {
		comment  string
		input    String
		reqLevel TrustLevel
		want     string
		wantErr  error
	}{
		{
			comment:  "valid url",
			input:    EscapeURL("http://google.com/"),
			reqLevel: URLSafe,
			want:     "http://google.com/",
			wantErr:  nil,
		},
		{
			comment:  "escaped url cannot be fully trusted",
			input:    EscapeURL("http://google.com/"),
			reqLevel: FullyTrusted,
			want:     "",
			wantErr:  ErrStringUntrusted,
		},
		{
			comment:  "constants can do whatever",
			input:    Const("Hello!"),
			reqLevel: FullyTrusted,
			want:     "Hello!",
			wantErr:  nil,
		},
		{
			comment:  "HTML is valid as text",
			input:    HTML{"<html></html>"},
			reqLevel: TextSafe,
			want:     "<html></html>",
			wantErr:  nil,
		},
	} {
		t.Run(tc.comment, func(t *testing.T) {
			s, err := Check(tc.input, tc.reqLevel)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Check(%v, %v) => (%q, %v), wanted error %v", tc.input, tc.reqLevel, s, err, tc.wantErr)
			}

			if s != tc.want {
				t.Errorf("Check(%v, %v) => (%q, %v), wanted string %q", tc.input, tc.reqLevel, s, err, tc.want)
			}
		})
	}
}
