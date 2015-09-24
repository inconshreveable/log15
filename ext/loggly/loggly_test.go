package loggly

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
	"unsafe"

	"gopkg.in/inconshreveable/log15.v2"
)

func TestLogglyHandler(t *testing.T) {

	const testPort = `:8123`

	type resultStructure struct {
		Context   map[string]interface{} `json:"context"`
		Hostname  string                 `json:"hostname"`
		Level     string                 `json:"level"`
		Message   string                 `json:"message"`
		Timestamp time.Time              `json:"timestamp"`
	}

	var testCases = []struct {
		Level   log15.Lvl
		Message string
		Context log15.Ctx
		Tags    []string
	}{
		{log15.LvlInfo, "a test message", log15.Ctx{"foo": "bar"}, []string{"foo", "bar"}},
		{log15.LvlWarn, "another test message", log15.Ctx{"foo": "bar", "number": 42}, nil},
	}

	for i, testCase := range testCases {
		// setup a listener that we can close
		listener, err := net.Listen("tcp", testPort)
		if err != nil {
			t.Fatal(err)
		}
		serverDoneCh := make(chan bool, 1)
		server := http.Server{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer r.Body.Close()
				fmt.Printf("%d server got request\n", i)

				var fail bool
				defer func(serverDoneCh chan bool) {
					serverDoneCh <- fail
					fmt.Printf("%d request fail is %t and ch pointer is: %x\n", i, fail, uintptr(unsafe.Pointer(&serverDoneCh)))
				}(serverDoneCh)

				result := &resultStructure{}
				err := json.NewDecoder(r.Body).Decode(result)
				if err != nil {
					t.Log(err)
					fail = true
					return
				}
				if result.Message != testCase.Message {
					t.Logf("Server got wrong log message. Expected: `%s`. Got: `%s`", testCase.Message, result.Message)
					fail = true
				}
				if result.Level != testCase.Level.String() {
					t.Logf("Server got wrong log level. Expected: `%s`. Got: `%s`", testCase.Level.String(), result.Level)
					fail = true
				}
				if testCase.Tags != nil {
					expectedTags := strings.Join(testCase.Tags, ",")
					gotTags := r.Header.Get("X-Loggly-Tag")
					if gotTags != expectedTags {
						t.Logf("Server got wrong tags. Expected `%s`. Got: `%s`", expectedTags, gotTags)
						fail = true
					}
				}
			}),
		}
		go func() {
			server.Serve(listener)
		}()

		time.Sleep(1 * time.Second)

		logger := log15.New()
		logglyHandler := NewLogglyHandler("test-token")
		if testCase.Tags != nil {
			logglyHandler.Tags = testCase.Tags
		}
		expectedEndpoint := `https://logs-01.loggly.com/inputs/test-token`
		if logglyHandler.Endpoint != expectedEndpoint {
			t.Errorf("invalid loggly endpoint. Expected `%s`. Got `%s`", expectedEndpoint, logglyHandler.Endpoint)
			return
		}
		logglyHandler.Endpoint = `http://localhost` + testPort + `/`
		logger.SetHandler(logglyHandler)
		switch testCase.Level {
		case log15.LvlInfo:
			logger.Info(testCase.Message, testCase.Context)
		case log15.LvlWarn:
			logger.Warn(testCase.Message, testCase.Context)
		}

		// Wait until request has been passed to server
		fmt.Printf("%d waiting for response on channel with pointer %x\n", i, uintptr(unsafe.Pointer(&serverDoneCh)))
		fail := <-serverDoneCh
		fmt.Printf("%d received fail is %t from channel with pointer %x\n", i, fail, uintptr(unsafe.Pointer(&serverDoneCh)))
		listener.Close()
		if fail {
			t.FailNow()
		}
		time.Sleep(1 * time.Second)
		fmt.Printf("%d done\n", i)
	}

	fmt.Println("am here")
}
