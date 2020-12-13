package html

import (
	"bufio"
	"bytes"
	"io"
	"unicode/utf8"
)

func fprintRawNewline(w io.Writer, depth int, indent string) error {
	if _, err := w.Write([]byte{'\n'}); err != nil {
		return err
	}
	for i := 0; i < depth; i++ {
		if _, err := io.WriteString(w, indent); err != nil {
			return err
		}
	}
	return nil
}

func fprintBlockText(w io.Writer, depth, width int, indent string, text io.Reader) error {
	if width <= 0 {
		_, err := io.Copy(w, text)
		return err
	}

	var line bytes.Buffer
	runeCount := depth * len(indent)

	scanner := bufio.NewScanner(text)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Bytes()
		l := utf8.RuneCount(word)
		// Flush on overrun, but don't create empty lines (could happen if the
		// word is longer than the line length and there are no bytes buffered.)
		if runeCount+l > width && line.Len() > 0 {
			if _, err := line.WriteTo(w); err != nil {
				return err
			}
			line.Reset()
			runeCount = depth*len(indent) + l
			if err := fprintRawNewline(w, depth, indent); err != nil {
				return err
			}
		}
		line.Write(word)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if line.Len() > 0 {
		_, err := line.WriteTo(w)
		return err
	}

	return nil
}
