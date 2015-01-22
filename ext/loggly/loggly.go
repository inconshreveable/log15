package loggly

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/inconshreveable/log15.v2"
)

// LogglyHandler sends logs to Loggly.
// LogglyHandler should be created by NewLogglyHandler.
// Exported fields can be modified during setup, but should not be touched when the Handler is in use.
// LogglyHandler implements log15.Handler
type LogglyHandler struct {
	// Client can be modified or replaced with a custom http.Client
	Client *http.Client

	// Defaults contains key/value items that are added to every log message.
	// Extra values can be added during the log15 setup.
	//
	// NewLogglyHandler adds a single record: "hostname", with the return value from os.Hostname().
	// When os.Hostname() returns with an error, the key "hostname" is not set and this map will be empty.
	Defaults map[string]interface{}

	// Tags are sent to loggly with the log.
	Tags []string

	// Endpoint is set to the https URI where logs are sent
	Endpoint string
}

// NewLogglyHandler creates a new LogglyHandler instance
// Exported field on the LogglyHandler can modified before it is being used.
func NewLogglyHandler(token string) *LogglyHandler {
	lh := &LogglyHandler{
		Endpoint: `https://logs-01.loggly.com/inputs/` + token,

		Client: &http.Client{},

		Defaults: make(map[string]interface{}),
	}

	// if hostname is retrievable, set it as extra field
	if hostname, err := os.Hostname(); err == nil {
		lh.Defaults["hostname"] = hostname
	}

	return lh
}

// Log sends the given *log15.Record to loggly.
// Standard fields are:
//     - message, the record's message.
//     - level, the record's level as string.
//     - timestamp, the record's timestamp in UTC timezone truncated to microseconds.
//     - context, (optional) the context fields from the record.
// Extra fields are the configurable with the LogglyHandler.Defaults map
// By default this contains:
//     - hostname, the system hostname
func (lh *LogglyHandler) Log(r *log15.Record) error {
	// create message structure
	msg := lh.createMessage(r)

	// send message
	err := lh.sendSingle(msg)
	if err != nil {
		return err
	}

	return nil
}

// createMessage takes a log15.Record and returns a loggly message structure
func (lh *LogglyHandler) createMessage(r *log15.Record) map[string]interface{} {
	// set standard values
	msg := map[string]interface{}{
		"message": r.Msg,
		"level":   r.Lvl.String(),
		// for loggly we need to truncate the timestamp to microsecond precision and convert it to UTC timezone
		"timestamp": r.Time.Truncate(time.Microsecond).In(time.UTC),
	}

	// apply defaults
	for key, value := range lh.Defaults {
		msg[key] = value
	}

	// optionally add context
	if len(r.Ctx) > 0 {
		context := make(map[string]interface{}, len(r.Ctx)/2)
		for i := 0; i < len(r.Ctx); i += 2 {
			key := r.Ctx[i]
			value := r.Ctx[i+1]
			keyStr, ok := key.(string)
			if !ok {
				keyStr = fmt.Sprintf("%v", key)
			}
			context[keyStr] = value
		}
		msg["context"] = context
	}

	// got a nice message to deliver
	return msg
}

// sendSingle sends a single loggly structure to their http endpoint
func (lh *LogglyHandler) sendSingle(msg map[string]interface{}) error {
	// encode the message to json
	postBuffer := &bytes.Buffer{}
	err := json.NewEncoder(postBuffer).Encode(msg)
	if err != nil {
		return err
	}

	// create request
	req, err := http.NewRequest("POST", lh.Endpoint, postBuffer)
	req.Header.Add("User-Agent", "log15")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(postBuffer.Len()))

	// apply tags
	if len(lh.Tags) > 0 {
		req.Header.Add("X-Loggly-Tag", strings.Join(lh.Tags, ","))
	}

	// do request
	resp, err := lh.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// check statuscode
	if resp.StatusCode != 200 {
		resp, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error: %s", string(resp))
	}

	// validate response
	response := &logglyResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}
	if response.Response != "ok" {
		return errors.New(`loggly response was not "ok"`)
	}

	// all done
	return nil
}
