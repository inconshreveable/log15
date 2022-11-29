package log15

type LegacyHandler interface {
	Log(r *Record) error
}

type LegacyLogger interface {
	New(ctx ...interface{}) LegacyLogger

	GetHandler() LegacyHandler

	SetHandler(h LegacyHandler)

	Debug(msg string, ctx ...interface{})
	Info(msg string, ctx ...interface{})
	Warn(msg string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
	Crit(msg string, ctx ...interface{})
}

type compat struct {
	Handler
	Logger
}

func (c *compat) Log(r *Record) error {
	return c.Handler.Log(*r)
}

func (c *compat) GetHandler() LegacyHandler {
	return &compat{Handler: c.Logger.GetHandler()}
}

func (c *compat) SetHandler(h LegacyHandler) {
	c.Logger.SetHandler(FuncHandler(func(r Record) error {
		return h.Log(&r)
	}))
}

func (c *compat) New(args ...interface{}) LegacyLogger {
	return &compat{Logger: c.Logger.New(args...)}
}

// CompatHandler wraps a handler for use with pre-v3 log15.Handler consumers.
func CompatHandler(h Handler) LegacyHandler {
	return &compat{Handler: h}
}

// CompatHandler wraps a handler for use with pre-v3 log15.Handler consumers.
func CompatLogger(l Logger) LegacyLogger {
	return &compat{Logger: l}
}
