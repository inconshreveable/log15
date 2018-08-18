package config

import (
	"fmt"
	"strings"

	"github.com/inconshreveable/log15"
)

type Fmt int

// GLEF is not a  format, it's a handler!!
const (
	FmtTerminal Fmt = iota
	FmtJson
	FmtLogfmt
)

// Returns the appropriate Lvl from a string name.
// Useful for parsing command line args and configuration files.
func FmtFromString(fmtString string) (Fmt, error) {
	switch strings.ToLower(fmtString) {
	case "terminal", "term", "console":
		return FmtTerminal, nil
	case "json":
		return FmtJson, nil
	case "logfmt":
		return FmtLogfmt, nil
	default:
		return FmtTerminal, fmt.Errorf("Unknown format: %v", fmtString)
	}
}

// UnmarshalString to implement StringUnmarshaller
func (f Fmt) UnmarshalString(from string) (interface{}, error) {
	f1, err := FmtFromString(from)
	if err != nil {
		return nil, err
	}

	return f1, nil
}

func (f Fmt) NewFormat() log15.Format {
	switch f {
	case FmtTerminal:
		return log15.TerminalFormat()
	case FmtJson:
		return log15.JsonFormat()
	case FmtLogfmt:
		return log15.LogfmtFormat()
	default:
		panic("invalid format: ")
	}
}
