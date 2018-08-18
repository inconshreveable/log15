package config

import (
	"fmt"
	"os"
	"testing"
	"time"

	"bufio"
	"encoding/json"
	"io"
	"path/filepath"

	"net"

	"github.com/inconshreveable/log15"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func testConfigLogger(conf string) (log15.Logger, error) {
	configMap := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(conf), &configMap)
	if err != nil {
		return nil, fmt.Errorf("yaml umnarshall failed: %v", err)
	}
	return Logger(configMap)
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

func testPrepareForFile(path string) error {
	lfile := "./testdata/temp/logTestLevelConfig.log"

	err := os.MkdirAll(filepath.Dir(lfile), 0777)
	if err != nil {
		return err
	}
	return os.Remove(lfile)
}

// givenHostAvailable skips if host not available
func givenHostAvailable(host string, t *testing.T) {
	ips, err := net.LookupIP(host)
	if err != nil || len(ips) == 0 {
		t.Skipf("can't resolve host '%v'", host)
	}
}

func TestReadSimpleConfig(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	var config = `
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

	l, err := testConfigLogger(config)
	require.Nil(err)

	l.Info("Hello, logs!")
	l.Debug("Hello, debug logs!")
}

func TestMultiConfig(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	var config = `
  level: INFO
  handlers:
    - kind: multi
      handlers:
        - kind: stdout
          format: terminal
        - kind: stderr
          format: json
        - kind: stdout
          format: logfmt
`

	l, err := testConfigLogger(config)
	require.Nil(err)

	l.Info("Hello, logs!")
	l.Debug("Hello, debug logs!")
}

func TestFailoverConfig(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	var config = `
  level: INFO
  handlers:
    - kind: failover
      handlers:
        - kind: stdout
          format: terminal
        - kind: stderr
          format: json
        - kind: stdout
          format: logfmt
`

	l, err := testConfigLogger(config)
	require.Nil(err)

	l.Info("Hello, logs!")
	l.Debug("Hello, debug logs!")
}

func TestMatchFilterConfig(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	var config = `
  level: info
  handlers:
    - kind: filter
      key: matcher
      value: foo
      handler:
        kind: stdout
        format: json
`

	l, err := testConfigLogger(config)
	require.Nil(err)

	l.Info("log foo", "matcher", "foo")
	l.Info("log bar", "matcher", "bar")
}

func TestGelfConfig(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	var config = `
  level: INFO
  handlers:
    - kind: gelf
      address: "logger:12201"
`

	givenHostAvailable("logger", t)

	l, err := testConfigLogger(config)
	require.Nil(err)

	l.Info("Hello, gelf!")
}

func TestLevelConfig(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var config = `
  level: Warn
  handlers:
    - kind: file
      path: ./testdata/temp/logTestLevelConfig.log
      format: json
      level: info
`
	lfile := "./testdata/temp/logTestLevelConfig.log"

	testPrepareForFile(lfile)

	l, err := testConfigLogger(config)
	require.Nil(err)

	l.Info("Hello, logs!")
	l.Debug("Hello, debug logs!", "mark", 1)

	r, err := testParseFile(lfile)
	require.Nil(err)
	//outputs.Dump(r, "records")
	require.EqualValues(1, len(r))

	assert.Equal("Hello, logs!", r[0]["msg"])
	//assert.EqualValues(1,r[1]["mark"])
}

// readJson is used to parse messages on the fly..
func readJson(rd *bufio.Reader) (map[string]interface{}, error) {
	s, err := rd.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("Failed to read string: %v", err)
	}
	var r map[string]interface{}
	json.Unmarshal(s, &r)
	return r, nil
}

func TestBufferedTcpNetHandler(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	require := require.New(t)

	listen, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", listen)
	}

	config := fmt.Sprintf(`
  level: Info
  handlers:
    - kind: buffer
      level: debug # w/o this, the nested handler(s) won't be activated!!
      handler:
        kind: net
        url: %v://%v
        format: json
        level: debug
`, listen.Addr().Network(), listen.Addr().String())

	//fmt.Println("config:", config)
	errs := make(chan error)
	go func() {
		c, err := listen.Accept()
		if err != nil {
			t.Fatalf("Failed to accept connection: %v", err)
		}
		rd := bufio.NewReader(c)

		r, err := readJson(rd)
		assert.EqualValues("Hello, logs!", r["msg"])
		assert.EqualValues(1, r["mark"])

		r, err = readJson(rd)
		assert.EqualValues("Hello, debug logs!", r["msg"])
		assert.EqualValues(2, r["mark"])

		errs <- nil
	}()

	l, err := testConfigLogger(config)
	require.Nil(err)

	l.Info("Hello, logs!", "mark", 1)
	l.Debug("Hello, debug logs!", "mark", 2)

	select {
	case <-time.After(time.Second):
		t.Fatalf("Test timed out!")
	case <-errs:
		// ok
	}
}

func serveUDP(addr string) (*net.UDPConn, error) {

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("ResolveUDPAddr('%s'): %s", addr, err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, fmt.Errorf("ListenUDP: %s", err)
	}

	return conn, nil
}

func readUDPJson(conn *net.UDPConn) (map[string]interface{}, error) {
	var r map[string]interface{}

	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFrom(buffer)
	if err != nil {
		return r, err
	}

	json.Unmarshal(buffer[:n], &r)
	return r, nil
}

func TestUdpNetHandler(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	require := require.New(t)

	conn, err := serveUDP("localhost:0")
	require.Nil(err)

	fmt.Println("r:", conn.LocalAddr())

	config := fmt.Sprintf(`
  level: Info
  handlers:
    - kind: net
      url: udp://%v
      format: json
      level: debug
`, conn.LocalAddr())

	//fmt.Println("config:", config)
	errs := make(chan error)
	go func() {
		r, err := readUDPJson(conn)
		require.Nil(err)

		assert.EqualValues("Hello, logs!", r["msg"])
		assert.EqualValues(1, r["mark"])

		r, err = readUDPJson(conn)
		require.Nil(err)

		assert.EqualValues("Hello, debug logs!", r["msg"])
		assert.EqualValues(2, r["mark"])

		errs <- nil
	}()

	l, err := testConfigLogger(config)
	require.Nil(err)

	l.Info("Hello, logs!", "mark", 1)
	l.Debug("Hello, debug logs!", "mark", 2)

	select {
	case <-time.After(time.Second):
		t.Fatalf("Test timed out!")
	case <-errs:
		// ok
	}
}
