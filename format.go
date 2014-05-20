package log15

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	timeLayout  = "2006-01-02T15:04:05-0700"
	floatFormat = 'f'
)

type Format interface {
	Format(r *Record) []byte
}

func FormatFunc(f func(*Record) []byte) Format {
	return formatFunc(f)
}

type formatFunc func(*Record) []byte

func (f formatFunc) Format(r *Record) []byte {
	return f(r)
}

// TerminalFormat formats log records optimized for human readability on
// a terminal with color-coded level output and terser human friendly timestamp.
// This format should only be used for interactive programs or while developing.
//
//     [TIME] [LEVEL] MESAGE key=value key=value ...
//
// Example:
//
//     [May 16 20:58:45] [DBUG] remove route ns=haproxy addr=127.0.0.1:50002
//
func TerminalFormat() Format {
	return FormatFunc(func(r *Record) []byte {
		var color = 0
		switch r.Lvl {
		case LvlCrit:
			color = 35
		case LvlError:
			color = 31
		case LvlWarn:
			color = 33
		case LvlInfo:
			color = 32
		case LvlDebug:
			color = 36
		}

		var s string
		lvl := strings.ToUpper(r.Lvl.String())
		if color > 0 {
			s = fmt.Sprintf("[%s] [\x1b[%dm%s\x1b[0m] %s ", r.Time.Format(time.Stamp), color, lvl, r.Msg)
		} else {
			s = fmt.Sprintf("[%s] [%s] %s ", r.Time.Format(time.Stamp), lvl, r.Msg)
		}

		bytes := logfmt(r.Ctx)
		return append([]byte(s), bytes...)
	})
}

// LogfmtFormat prints records in logfmt format, an easy machine-parseable but human-readable
// format for key/value pairs.
//
// For more details see: http://godoc.org/github.com/kr/logfmt
//
func LogfmtFormat() Format {
	return FormatFunc(func(r *Record) []byte {
		common := []interface{}{"t", r.Time, "lvl", r.Lvl, "msg", r.Msg}
		return logfmt(append(common, r.Ctx...))
	})
}

func logfmt(ctx []interface{}) []byte {
	pieces := make([]string, 0)

	for i := 0; i < len(ctx); i += 2 {
		k, ok := ctx[i].(string)
		var s string
		if !ok {
			s = fmt.Sprintf(`%s="%+v is not a string key"`, errorKey, ctx[i])
		} else {
			// XXX: we should probably check that all of your key bytes aren't invalid`
			s = fmt.Sprintf(`%s=%s`, k, formatLogfmtValue(ctx[i+1]))
		}

		pieces = append(pieces, s)
	}

	return []byte(strings.Join(append(pieces, "\n"), " "))
}

// JsonFormat formats log records as JSON objects separated by newlines.
// It is the equivalent of JsonFormatEx(false, true).
func JsonFormat() Format {
	return JsonFormatEx(false, true)
}

// JsonFormatEx formats log records as JSON objects. If pretty is true,
// records will be pretty-printed. If lineSeparated is true, records
// will be logged with a new line between each record.
func JsonFormatEx(pretty, lineSeparated bool) Format {
	jsonMarshal := json.Marshal
	if pretty {
		jsonMarshal = func(v interface{}) ([]byte, error) {
			return json.MarshalIndent(v, "", "    ")
		}
	}

	return FormatFunc(func(r *Record) []byte {
		props := make(map[string]interface{})

		props["t"] = r.Time
		props["lvl"] = r.Lvl
		props["msg"] = r.Msg

		for i := 0; i < len(r.Ctx); i += 2 {
			k, ok := r.Ctx[i].(string)
			if !ok {
				props[errorKey] = fmt.Sprintf("%+v is not a string key", r.Ctx[i])
			}
			props[k] = formatJsonValue(r.Ctx[i+1])
		}

		b, err := jsonMarshal(props)
		if err != nil {
			b, _ = jsonMarshal(map[string]string{
				errorKey: err.Error(),
			})
			return b
		}

		if lineSeparated {
			b = append(b, '\n')
		}

		return b
	})
}

func formatShared(value interface{}) interface{} {
	switch v := value.(type) {
	case time.Time:
		return v.Format(timeLayout)

	case error:
		return v.Error()

	case fmt.Stringer:
		return v.String()

	default:
		return v
	}
}

func formatJsonValue(value interface{}) interface{} {
	value = formatShared(value)
	switch value.(type) {
	case int, int8, int16, int32, int64, float32, float64, uint, uint8, uint16, uint32, uint64, string:
		return value
	default:
		return fmt.Sprintf("%+v", value)
	}
}

// formatValue formats a value for serialization
func formatLogfmtValue(value interface{}) string {
	if value == nil {
		return "nil"
	}

	value = formatShared(value)
	switch v := value.(type) {
	case string:
		return escapeString(v)

	case bool:
		return strconv.FormatBool(v)

	case int:
		return strconv.FormatInt(int64(v), 10)

	case int8:
		return strconv.FormatInt(int64(v), 10)

	case int16:
		return strconv.FormatInt(int64(v), 10)

	case int32:
		return strconv.FormatInt(int64(v), 10)

	case int64:
		return strconv.FormatInt(v, 10)

	case float32:
		return strconv.FormatFloat(float64(v), floatFormat, 3, 64)

	case float64:
		return strconv.FormatFloat(v, floatFormat, 3, 64)

	case uint:
		return strconv.FormatUint(uint64(v), 10)

	case uint8:
		return strconv.FormatUint(uint64(v), 10)

	case uint16:
		return strconv.FormatUint(uint64(v), 10)

	case uint32:
		return strconv.FormatUint(uint64(v), 10)

	case uint64:
		return strconv.FormatUint(v, 10)

	default:
		return escapeString(fmt.Sprintf("%+v", value))
	}
}

func escapeString(s string) string {
	needQuotes := false
	e := new(bytes.Buffer)
	e.WriteByte('"')
	for _, r := range s {
		if r <= ' ' || r == '=' || r == '"' {
			needQuotes = true
		}

		switch r {
		case '\\', '"':
			e.WriteByte('\\')
			e.WriteByte(byte(r))
		case '\n':
			e.WriteByte('\\')
			e.WriteByte('n')
		case '\r':
			e.WriteByte('\\')
			e.WriteByte('r')
		case '\t':
			e.WriteByte('\\')
			e.WriteByte('t')
		default:
			e.WriteRune(r)
		}
	}
	e.WriteByte('"')
	start, stop := 0, e.Len()
	if !needQuotes {
		start, stop = 1, stop-1
	}
	return string(e.Bytes()[start:stop])
}
