package gelf

import (
	"testing"
	"time"
)

func TestCtxToMap(t *testing.T) {

	loc, err := time.LoadLocation("Europe/Vienna")
	if err != nil {
		t.Fatalf("can't load Timezone: %v", err)
	}
	logTime := time.Date(2016, 11, 23, 13, 01, 02, 123100*1e3, loc)

	expected := map[string]interface{}{
		"_msg":    "a message",
		"_foo":    "baz",
		"_number": 1,
		"_t":      logTime,
	}
	ctx := []interface{}{"msg", "a message", "foo", "bar", "foo", "baz", "number", 1, "t", logTime}

	cm := CtxToMap(ctx)

	for k, v := range expected {
		if cm[k] != v {
			t.Fatalf("%v: expected: '%v', got: %v", k, v, cm[k])
		}
	}
}
