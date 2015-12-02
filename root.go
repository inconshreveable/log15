package log15

import (
	"fmt"
	"os"

	"github.com/inconshreveable/log15/term"
	"github.com/mattn/go-colorable"
)

var (
	root          *logger
	StdoutHandler = StreamHandler(os.Stdout, LogfmtFormat())
	StderrHandler = StreamHandler(os.Stderr, LogfmtFormat())
)

func init() {
	if term.IsTty(os.Stdout.Fd()) {
		StdoutHandler = StreamHandler(colorable.NewColorableStdout(), TerminalFormat())
	}

	if term.IsTty(os.Stderr.Fd()) {
		StderrHandler = StreamHandler(colorable.NewColorableStderr(), TerminalFormat())
	}

	root = &logger{[]interface{}{}, new(swapHandler)}
	root.SetHandler(StdoutHandler)
}

// New returns a new logger with the given context.
// New is a convenient alias for Root().New
func New(ctx ...interface{}) Logger {
	return root.New(ctx...)
}

// Root returns the root logger
func Root() Logger {
	return root
}

// The following functions bypass the exported logger methods (logger.Debug,
// etc.) to keep the call depth the same for all paths to logger.write so
// runtime.Caller(2) always refers to the call site in client code.

// Debug is a convenient alias for Root().Debug
func Debug(msg interface{}, ctx ...interface{}) {
	root.write(fmt.Sprint(msg), LvlDebug, ctx)
}

// Info is a convenient alias for Root().Info
func Info(msg interface{}, ctx ...interface{}) {
	root.write(fmt.Sprint(msg), LvlInfo, ctx)
}

// Warn is a convenient alias for Root().Warn
func Warn(msg interface{}, ctx ...interface{}) {
	root.write(fmt.Sprint(msg), LvlWarn, ctx)
}

// Error is a convenient alias for Root().Error
func Error(msg interface{}, ctx ...interface{}) {
	root.write(fmt.Sprint(msg), LvlError, ctx)
}

// Crit is a convenient alias for Root().Crit
func Crit(msg interface{}, ctx ...interface{}) {
	root.write(fmt.Sprint(msg), LvlCrit, ctx)
}
