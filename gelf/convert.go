package gelf

import (
	"strconv"
	"strings"
)

func CtxToMap(ctx []interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(ctx)/2)
	for i := 0; i < len(ctx); i += 2 {
		s := ctx[i].(string)
		m["_"+s] = ctx[i+1]
	}

	return m
}

func ShortAndFull(msg string) (short string, full string) {
	lines := strings.SplitN(msg, "\n", 2)
	short = lines[0]
	if len(lines) > 1 {
		full = msg
	}

	return short, full
}

// caller searches a context list for an entry called "caller" and splits it into
// filename and line number.
func Caller(ctx map[string]interface{}) (string, int) {
	info, present := ctx["_caller"]
	if !present {
		return "", 0
	}

	parts := strings.Split(info.(string), ":")
	line, _ := strconv.Atoi(parts[1])
	return parts[0], line
}
