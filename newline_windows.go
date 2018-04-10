// +build windows

package log15

import "bytes"

func newLine(buf *bytes.Buffer) {
	buf.WriteString("\r\n")
}

func newLineJson(b []byte) []byte {
	return append(b, '\r', '\n')
}
