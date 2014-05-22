package ext

import log "gopkg.in/inconshreveable/log15.v1"

// ErrorHandler wraps another handler and passes all records through
// unchanged except if the logged context contains a non-nil error
// value in its context. Then ErrorHandler will *increase* the log
// level to LvlError, unless it was already at LvlCrit.
//
// This allows you to log the result of all functions for debugging
// and capture error conditions when in production with a single
// log line. Example:
//
//     reply, err := redisConn.Do("SET", "foo", "bar")
//     logger.Debug("Wrote value to redis", "reply", reply, "err", err)
//     if err != nil {
//         return err
//     }
//
func ErrorHandler(h log.Handler) log.Handler {
	return errorHandler{h}
}

type errorHandler struct {
	h log.Handler
}

func (h errorHandler) Log(r *log.Record) error {
	if r.Lvl < log.LvlError {
		for i := 1; i < len(r.Ctx); i++ {
			if v, ok := r.Ctx[i].(error); ok && v != nil {
				r.Lvl = log.LvlError
				break
			}
		}
	}

	return h.Log(r)
}
