package term

import "syscall"

const ioctlReadTermios = syscall.TIOCGETA

// Termios functions describe a general terminal interface that is
// provided to control asynchronous communications ports.
//
// Deprecated: this type is kept for source compatibility only.
type Termios syscall.Termios
