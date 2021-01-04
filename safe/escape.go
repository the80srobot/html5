package safe

import (
	"errors"
	"fmt"
	"html"
	"net/url"
)

var (
	errInvalidInput    = errors.New("invalid string input")
	errForbiddenSchema = errors.New("URL schema not allowed")
)

func escapeHTML(s string) (string, error) {
	return "", errors.New("HTML fragments must currently be fully trusted")
}

func escapeText(s string) (string, error) {
	return html.EscapeString(s), nil
}

func escapeURL(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", fmt.Errorf("%w: %v", errInvalidInput, err)
	}

	switch u.Scheme {
	case "http", "https", "mailto", "ftp", "":
		return u.String(), nil
	default:
		return "", fmt.Errorf("%w: %s", errForbiddenSchema, u.Scheme)
	}
}

func escapeAttribute(s string) (string, error) {
	return "", errors.New("attributes must currently be fully trusted")
}
