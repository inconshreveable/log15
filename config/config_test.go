package config

import (
	"fmt"
	"os"
	"testing"
	//	_ "geger.at/logsExplorer/log15/graylog"
	"bufio"
	"encoding/json"
	"github.com/inconshreveable/log15"
	"io"
	"path/filepath"
	"github.com/gernoteger/mapstructure-hooks"
	"gopkg.in/yaml.v2"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)


func testConfigLogger(config string) (log15.Logger, error) {
	c := LoggerConfig{}

	ci := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(config), &ci)
	if err != nil {
		return nil, err
	}

	// das gehört zusammen!!
	err = hooks.Decode(ci, &c)
	if err != nil {
		return nil, err
	}
	l, err := c.NewLogger()

	return l, err
}

func testParseFile(path string) ([]map[string]interface{}, error) {

	if file, err := os.Open(path); err == nil {

		// make sure it gets closed
		defer file.Close()
		return testParseReader(file)
	} else {
		return nil, err
	}
}

// parse to records until scanner closes
func testParseReader(file io.Reader) ([]map[string]interface{}, error) {
	records := make([]map[string]interface{}, 0)

	// create a new scanner and read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		m := make(map[string]interface{})
		err := json.Unmarshal(scanner.Bytes(), &m)
		if err != nil {
			return records, err
		}
		records = append(records, m) // not very efficient, but..
	}

	// check for errors
	if err := scanner.Err(); err != nil {
		return records, err
	}

	return records, nil
}

func testPrepareForFile(path string) error{
	lfile := "./testdata/temp/logTestLevelConfig.log"

	err:=os.MkdirAll(filepath.Dir(lfile), 0777)
	if err != nil {
		return err
	}
	return os.Remove(lfile)
}


func TestReadSimpleConfig(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	var t1 = `
  level: INFO
  handlers:
    - kind: stdout
      format: terminal
      level: info
    - kind: stderr
      format: json
      level: debug
    - kind: stdout
      format: logfmt
      level: info
`

	l, err := testConfigLogger(t1)
	require.Nil(err)

	l.Info("Hello, logs!")
	l.Debug("Hello, debug logs!")
}

func TestLevelConfig(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var t1 = `
  level: Warn
  handlers:
    - kind: file
      path: ./testdata/temp/logTestLevelConfig.log
      format: json
      level: info
`
	lfile := "./testdata/temp/logTestLevelConfig.log"

	testPrepareForFile(lfile)

	l, err := testConfigLogger(t1)
	require.Nil(err)

	l.Info("Hello, logs!")
	l.Debug("Hello, debug logs!","mark",1)

	r, err := testParseFile(lfile)
	require.Nil(err)
	//outputs.Dump(r, "records")
	require.EqualValues(1,len(r))

	assert.Equal("Hello, logs!",r[0]["msg"])
	//assert.EqualValues(1,r[1]["mark"])
}
