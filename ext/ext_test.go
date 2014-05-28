package ext

import (
	log "github.com/inconshreveable/log15"
	"testing"
)

type testHandler struct {
    r log.Record
}
func (h *testHandler) Log(r *log.Record) error {
    h.r = *r
    return nil
}

func TestHotSwapHandler(t *testing.T) {
	h1 := &testHandler{}

	l := log.New()
	h := HotSwapHandler(h1)
	l.SetHandler(h)

	l.Info("to h1")
	if h1.r.Msg != "to h1" {
		t.Fatalf("didn't get expected message to h1")
	}

	h2 := &testHandler{}
	h.Swap(h2)
	l.Info("to h2")
	if h2.r.Msg != "to h2" {
		t.Fatalf("didn't get expected message to h2")
	}
}
