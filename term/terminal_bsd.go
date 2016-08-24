// +build freebsd openbsd netbsd dragonfly darwin
package term

import "syscall"

const ioctlReadTermios = syscall.TIOCGETA

type Termios syscall.Termios
