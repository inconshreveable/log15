package ext

import (
	"flag"

	"github.com/inconshreveable/log15"
)

type lvlFlag struct {
	lvl *log15.Lvl
}

// FlagLvlValue creates an adapter that satifies the flag.Value interface.
func FlagLvlValue(lvl *log15.Lvl) flag.Value {
	return lvlFlag{lvl: lvl}
}

// Set tries to set the flag to the specified string value.
func (l lvlFlag) Set(lvlString string) error {
	lvl, err := log15.LvlFromString(lvlString)
	if err != nil {
		return err
	}
	*l.lvl = lvl
	return nil
}

// String returns a textual representation of the flag value.
func (l lvlFlag) String() string {
	return l.lvl.String()
}

// FlagLvl registers a flag on the default flag.CommandLine flag-set. If you
// need more control, use the FlagLvlValue adapter directly.
func FlagLvl(name string, value log15.Lvl, usage string) *log15.Lvl {
	flag.Var(FlagLvlValue(&value), name, usage)
	return &value
}

// FlagLvlVar registers a flag on the default flag.CommandLine flag-set. If
// you need more control, use the FlagLvlValue adapter directly.
func FlagLvlVar(p *log15.Lvl, name string, value log15.Lvl, usage string) {
	*p = value
	flag.Var(FlagLvlValue(p), name, usage)
}
