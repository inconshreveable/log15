// +build !windows
package log15

import "bytes"

func newLine(buf *bytes.Buffer) {
  buf.WriteByte('\n')
}
