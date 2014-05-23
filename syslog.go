// +build !windows,!plan9

package log15

import (
	"log/syslog"
	"strings"
)

type syslogHandler struct {
	sysWr *syslog.Writer
	fmtr  Format
}

func (h *syslogHandler) Log(r *Record) error {
	var syslogFn = h.sysWr.Info
	switch r.Lvl {
	case LvlCrit:
		syslogFn = h.sysWr.Crit
	case LvlError:
		syslogFn = h.sysWr.Err
	case LvlWarn:
		syslogFn = h.sysWr.Warning
	case LvlInfo:
		syslogFn = h.sysWr.Info
	case LvlDebug:
		syslogFn = h.sysWr.Debug
	}

	s := strings.TrimSpace(string(h.fmtr.Format(r)))
	return syslogFn(s)
}

// SyslogHandler opens a connection to the system syslog daemon by calling
// syslog.New and writes all records to it.
func SyslogHandler(priority syslog.Priority, tag string, fmtr Format) (Handler, error) {
	wr, err := syslog.New(priority, tag)
	return sharedSyslog(fmtr, wr, err)
}

// SyslogHandler opens a connection to a log daemon over the network and writes
// all log records to it.
func SyslogNetHandler(net, addr string, priority syslog.Priority, tag string, fmtr Format) (Handler, error) {
	wr, err := syslog.Dial(net, addr, priority, tag)
	return sharedSyslog(fmtr, wr, err)
}

func sharedSyslog(fmtr Format, sysWr *syslog.Writer, err error) (Handler, error) {
	if err != nil {
		return nil, err
	}

	return LazyHandler(&closingHandler{sysWr, &syslogHandler{sysWr, fmtr}}), nil
}

func (m muster) SyslogHandler(priority syslog.Priority, tag string, fmtr Format) Handler {
	return must(SyslogHandler(priority, tag, fmtr))
}

func (m muster) SyslogNetHandler(net, addr string, priority syslog.Priority, tag string, fmtr Format) Handler {
	return must(SyslogNetHandler(net, addr, priority, tag, fmtr))
}
