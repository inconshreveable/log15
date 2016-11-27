package gelfconfig

import (
	"testing"
	"gopkg.in/yaml.v2"
	"github.com/inconshreveable/log15"
	"github.com/inconshreveable/log15/config"
	"github.com/stretchr/testify/require"
)

func testConfigLogger(conf string) (log15.Logger, error) {
	configMap := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(conf), &configMap)
	if err != nil {
		return nil, err
	}
	return config.Logger(configMap)
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
    - kind: gelf
      address: "web:12201"
`

	l, err := testConfigLogger(t1)
	require.Nil(err)

	l.Info("Hello, logs!","mark","me")
	l.Debug("Hello, debug logs!")

	//time.Sleep(time.Millisecond*100)
}
