package graylog

import (
	"testing"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"github.com/inconshreveable/log15"

	"geger.at/logsExplorer/log15/config"
	"github.com/gernoteger/mapstructure-hooks"
)

func testConfigLogger(conf string) (log15.Logger, error) {
	c := config.LoggerConfig{}

	ci := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(conf), &ci)
	if err != nil {
		return nil, err
	}
	err = hooks.Decode(ci, &c)
	if err != nil {
		return nil, err
	}
	l, err := c.NewLogger()

	return l, err
}


func TestUseSimpleConfigWithGelf(t *testing.T) {
	//assert := assert.New(t)
	require := require.New(t)

	var t1 = `
  level: INFO
  extra:
      mark: test
      user: gernot

  handlers:
    - kind: stdout
      format: terminal
      level: info
    - kind: stderr
      format: json
      level: info
    - kind: stdout
      format: logfmt
      level: info
    - kind: gelf
      address: "web:12201"
`

	l, err := testConfigLogger(t1)
	require.Nil(err)

	l.Info("Hello, logs!","mark","me")
	l.Debug("Hello, debug logs!")

	//time.Sleep(time.Millisecond*100)
}
