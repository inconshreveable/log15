// +build !windows,!plan9

package config

import (
	"log/syslog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func must(v interface{}, err error) interface{} {

	if err != nil {
		panic(err)
	}
	return v
}

func TestFacility(t *testing.T) {
	assert := assert.New(t)

	data := []struct {
		n string
		p syslog.Priority
	}{
		{"kern", syslog.LOG_KERN},
		{"user", syslog.LOG_USER},
		{"Mail", syslog.LOG_MAIL},
		{"daemOn", syslog.LOG_DAEMON},
		{"auth", syslog.LOG_AUTH},
		{"syslog", syslog.LOG_SYSLOG},
		{"lpr", syslog.LOG_LPR},
		{"news", syslog.LOG_NEWS},
		{"uucp", syslog.LOG_UUCP},
		{"cron", syslog.LOG_CRON},
		{"authpriv", syslog.LOG_AUTHPRIV},
		{"ftp", syslog.LOG_FTP},
		{"local0", syslog.LOG_LOCAL0},
		{"local1", syslog.LOG_LOCAL1},
		{"local2", syslog.LOG_LOCAL2},
		{"local3", syslog.LOG_LOCAL3},
		{"local4", syslog.LOG_LOCAL4},
		{"local5", syslog.LOG_LOCAL5},
		{"local6", syslog.LOG_LOCAL6},
		{"local7", syslog.LOG_LOCAL7},
		//		{"local7", syslog.LOG_LOCAL7},
	}

	for _, tt := range data {
		p, err := priority(tt.n)
		assert.Nil(err, "for '%v'", tt.n)
		assert.EqualValues(tt.p, p, "for '%v'", tt.n)
	}

	_, err := priority("")
	assert.Error(err)
}

func TestSyslogConfig(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	var config = `
  level: debug
  handlers:
    - kind: syslog
      tag: testing
      facility: local6
`

	l, err := testConfigLogger(config)
	require.Nil(err)

	// check manually with e.g. journalctl -f; I don't want to code it, because this is very systemdependent
	l.Info("Hello, syslog!", "mark", "me")
	l.Debug("debug")
	l.Error("error")
	l.Crit("crit")
}

func TestSyslogNetConfig(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	var config = `
  level: INFO
  handlers:
    - kind: syslog
      tag: testingnet
      facility: local6
      #url: tcp://web:1514
      url: udp://web:1514
`

	l, err := testConfigLogger(config)
	require.Nil(err)

	// check manually with e.g. journalctl -f; I don't want to code it, because this is very system dependent
	// doesn't seem to work properly with graylog..
	l.Info("Hello, syslog!", "mark", "me")
}
