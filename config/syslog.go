// +build !windows,!plan9

package config

import (
	"fmt"
	"log/syslog"
	"net/url"
	"strings"

	"github.com/gernoteger/mapstructure-hooks"
	"github.com/inconshreveable/log15"
)

func init() {
	hooks.Register(HandlerConfigType, "syslog", NewSyslogConfig)
}

// see https://tools.ietf.org/html/rfc3164 and https://en.wikipedia.org/wiki/Syslog
var SyslogFacilities = map[string]syslog.Priority{
	"kern":     syslog.LOG_KERN,
	"user":     syslog.LOG_USER,
	"mail":     syslog.LOG_MAIL,
	"daemon":   syslog.LOG_DAEMON,
	"auth":     syslog.LOG_AUTH,
	"syslog":   syslog.LOG_SYSLOG,
	"lpr":      syslog.LOG_LPR,
	"news":     syslog.LOG_NEWS,
	"uucp":     syslog.LOG_UUCP,
	"cron":     syslog.LOG_CRON,
	"authpriv": syslog.LOG_AUTHPRIV,
	"ftp":      syslog.LOG_FTP,
	//"ntp":      syslog.               ,
	//"audit":    syslog.               ,
	//"alert":    syslog.               ,
	//"cron":     syslog.               ,
	"local0": syslog.LOG_LOCAL0,
	"local1": syslog.LOG_LOCAL1,
	"local2": syslog.LOG_LOCAL2,
	"local3": syslog.LOG_LOCAL3,
	"local4": syslog.LOG_LOCAL4,
	"local5": syslog.LOG_LOCAL5,
	"local6": syslog.LOG_LOCAL6,
	"local7": syslog.LOG_LOCAL7,
}

// get numerical value from string; either from map, or decode as int
func priority(facility string) (syslog.Priority, error) {
	p, ok := SyslogFacilities[strings.ToLower(facility)]

	if ok {
		return p, nil
	}

	return 0, fmt.Errorf("invalid facility: '%v' use on of %#v", facility, SyslogFacilities)
}

// BufferConfig is a buffered handkler
type SyslogConfig struct {
	LevelHandlerConfig `mapstructure:",squash"`
	Format             Fmt

	URL string // url; if omitted, local syslog is used

	// see https://en.wikipedia.org/wiki/Syslog
	Facility string // kern.. or 0..24, see SyslogFacilities
	// from message Severity string // emerg,alert,crit,err,warning,notice,info,debug or 0..7

	Tag string // typical name of application
}

// make sure its's the right interface
var _ HandlerConfig = (*SyslogConfig)(nil)

func NewSyslogConfig() interface{} {
	return &SyslogConfig{}
}

func (c *SyslogConfig) NewHandler() (log15.Handler, error) {

	if c.Tag == "" {
		return nil, fmt.Errorf("SyslogConfig: tag required")
	}
	p, err := priority(c.Facility)
	if err != nil {
		return nil, err
	}

	if c.URL != "" {
		// we have an url: use net handler
		u, err := url.Parse(c.URL)
		if err != nil {
			return nil, err
		}

		return log15.SyslogNetHandler(u.Scheme, u.Host, p, c.Tag, c.Format.NewFormat())
	}
	return log15.SyslogHandler(p, c.Tag, c.Format.NewFormat())
}
