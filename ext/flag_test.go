package ext

import (
	"flag"
	"fmt"
	"testing"

	"github.com/inconshreveable/log15"
)

func TestFlagLvlValue(t *testing.T) {
	lvl := log15.LvlInfo

	value := FlagLvlValue(&lvl)

	if value.String() != log15.LvlInfo.String() {
		t.Errorf("invalid string: expected %v, got %v", log15.LvlInfo, value)
	}

	if err := value.Set(log15.LvlError.String()); err != nil {
		t.Fatalf("failed to set value: %v", err)
	}

	if lvl != log15.LvlError {
		t.Fatalf("invalid value: expected %v, got %v", log15.LvlError, lvl)
	}
}

func ExampleFlagLvlValue() {
	lvl := log15.LvlInfo
	flag.Var(FlagLvlValue(&lvl), "loglevel", "maximum log level")
}

func ExampleFlagLvl() {
	lvl := FlagLvl("loglevel", log15.LvlInfo, "maximum log level")
	flag.Parse()

	fmt.Println("maximum log level:", *lvl)
}

func ExampleFlagLvlVar() {
	var lvl log15.Lvl

	FlagLvlVar(&lvl, "loglevel", log15.LvlInfo, "maximum log level")
	flag.Parse()

	fmt.Println("maximum log level:", lvl)
}
