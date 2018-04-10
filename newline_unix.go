// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package log15

import "bytes"

func newLine(buf *bytes.Buffer) {
	buf.WriteByte('\n')
}
