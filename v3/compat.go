package log15

import (
	legacy "github.com/inconshreveable/log15"
)

type compat struct {
	Handler
	Logger
}

func fromLegacyRecord(r *legacy.Record) Record {
	return Record{
		Time: r.Time,
		Lvl:  Lvl(r.Lvl),
		Msg:  r.Msg,
		Ctx:  r.Ctx,
		KeyNames: &RecordKeyNames{
			Lvl:  r.KeyNames.Lvl,
			Msg:  r.KeyNames.Msg,
			Time: r.KeyNames.Time,
		},
	}
}

func toLegacyRecord(r Record) *legacy.Record {
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
	c.Logger.SetHandler(FuncHandler(func(r Record) error {
		return h.Log(toLegacyRecord(r))
	}))
}

func (c *compat) New(args ...interface{}) legacy.Logger {
	return &compat{Logger: c.Logger.New(args...)}
}

// CompatHandler wraps a handler for use with pre-v3 log15.Handler consumers.
func CompatHandler(h Handler) legacy.Handler {
	return &compat{Handler: h}
}

// CompatLogger wraps a Logger for use with pre-v3 log15.Logger consumers.
func CompatLogger(l Logger) legacy.Logger {
	return &compat{Logger: l}
}

// Interface assertions to make sure we got everything right
var _ interface {
	legacy.Logger
	legacy.Handler
} = (*compat)(nil)
