package ext

import (
	log "gopkg.in/inconshreveable/log15.v1"
	"sync"
)

// ErrorHandler wraps another handler and passes all records through
// unchanged except if the logged context contains a non-nil error
// value in its context. Then ErrorHandler will *increase* the log
// level to LvlError, unless it was already at LvlCrit.
//
// This allows you to log the result of all functions for debugging
// and capture error conditions when in production with a single
// log line. Example:
//
//     reply, err := redisConn.Do("SET", "foo", "bar")
//     logger.Debug("Wrote value to redis", "reply", reply, "err", err)
//     if err != nil {
//         return err
//     }
//
func ErrorHandler(h log.Handler) log.Handler {
	return errorHandler{h}
}

type errorHandler struct {
	h log.Handler
}

func (h errorHandler) Log(r *log.Record) error {
	if r.Lvl < log.LvlError {
		for i := 1; i < len(r.Ctx); i++ {
			if v, ok := r.Ctx[i].(error); ok && v != nil {
				r.Lvl = log.LvlError
				break
			}
		}
	}

	return h.Log(r)
}

// SpeculativeHandler is a handler for speculative logging. It
// keeps a ring buffer of the given size full of the last events
// logged into it. When Flush is called, all buffered log records
// are written to the wrapped handler. This is extremely for
// continuosly capturing debug level output, but only flushing it
// if an exception condition is encountered, or the user requests
// it.
func SpeculativeHandler(size int, h log.Handler) *Speculative {
	return &Speculative{
		handler: h,
		recs:    make([]*log.Record, size),
	}
}

type Speculative struct {
	mu      sync.Mutex
	idx     int
	recs    []*log.Record
	handler log.Handler
	full    bool
}

func (h *Speculative) Log(r *log.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.recs[h.idx] = r
	h.idx = (h.idx + 1) % len(h.recs)
	h.full = h.idx == 0
	return nil
}

func (h *Speculative) Flush() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.full {
		for _, r := range h.recs[h.idx:] {
			h.handler.Log(r)
		}
	}

	for _, r := range h.recs[:h.idx] {
		h.handler.Log(r)
	}

	// reset state
	h.full = false
	h.idx = 0
}

// HotSwapHandler wraps another handler that may swapped out
// dynamically at runtime in a thread-safe fashion.
// HotSwapHandler is the same functionality
// used to implement the SetHandler method for the default
// implementation of Logger.
func HotSwapHandler(h log.Handler) (*HotSwap) {
	return &HotSwap{handler: h}
}

type HotSwap struct {
	mu      sync.RWMutex
	handler log.Handler
}

func (h *HotSwap) Log(r *log.Record) error {
	defer h.mu.RUnlock()
	h.mu.RLock()
	err := h.handler.Log(r)
	return err
}

func (h *HotSwap) Swap(newHandler log.Handler) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handler = newHandler
}
