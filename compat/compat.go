// This file provides backwards compatibility with pre-v3 log15 code.
//
// The compatibility layer allows existing code that uses the legacy log15 API
// to work seamlessly with the v3 version. It provides adapter types and functions
// that translate between the old and new interfaces, enabling incremental migration
// from legacy code to the new v3 API.
//
// Users may want this when:
//   - Migrating existing applications to use v3 gradually
//   - Working with third-party libraries that still use the legacy log15 API
//   - Maintaining backwards compatibility while adopting new v3 features
//
// Use CompatHandler() and CompatLogger() to wrap v3 types for legacy code.

package compat

import (
	legacy "github.com/inconshreveable/log15"
	"github.com/inconshreveable/log15/v3"
)

type compat struct {
	log15.Handler
	log15.Logger
}

func fromLegacyRecord(r *legacy.Record) log15.Record {
	return log15.Record{
		Time: r.Time,
		Lvl:  log15.Lvl(r.Lvl),
		Msg:  r.Msg,
		Ctx:  r.Ctx,
		KeyNames: &log15.RecordKeyNames{
			Lvl:  r.KeyNames.Lvl,
			Msg:  r.KeyNames.Msg,
			Time: r.KeyNames.Time,
		},
	}
}

func toLegacyRecord(r log15.Record) *legacy.Record {
	var keyNames legacy.RecordKeyNames
	if r.KeyNames != nil {
		keyNames.Lvl = r.KeyNames.Lvl
		keyNames.Msg = r.KeyNames.Msg
		keyNames.Time = r.KeyNames.Time
	}
	return &legacy.Record{
		Time:     r.Time,
		Lvl:      legacy.Lvl(r.Lvl),
		Msg:      r.Msg,
		Ctx:      r.Ctx,
		KeyNames: keyNames,
	}
}

func (c *compat) Log(r *legacy.Record) error {
	return c.Handler.Log(fromLegacyRecord(r))
}

func (c *compat) GetHandler() legacy.Handler {
	return &compat{Handler: c.Logger.GetHandler()}
}

func (c *compat) SetHandler(h legacy.Handler) {
	c.Logger.SetHandler(log15.FuncHandler(func(r log15.Record) error {
		return h.Log(toLegacyRecord(r))
	}))
}

func (c *compat) New(args ...interface{}) legacy.Logger {
	return &compat{Logger: c.Logger.New(args...)}
}

// CompatHandler wraps a handler for use with pre-v3 log15.Handler consumers.
func CompatHandler(h log15.Handler) legacy.Handler {
	return &compat{Handler: h}
}

// CompatLogger wraps a Logger for use with pre-v3 log15.Logger consumers.
func CompatLogger(l log15.Logger) legacy.Logger {
	return &compat{Logger: l}
}

// Interface assertions to make sure we got everything right
var _ interface {
	legacy.Logger
	legacy.Handler
} = (*compat)(nil)
