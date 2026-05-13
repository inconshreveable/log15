package term

import "golang.org/x/term"

// IsTty returns true if the given file descriptor is a terminal.
//
// Deprecated: use golang.org/x/term.IsTerminal instead.
func IsTty(fd uintptr) bool {
	return term.IsTerminal(int(fd))
}
