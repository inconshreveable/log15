package loggly

//go:generate ffjson $GOFILE

// logglyResponse defines the json returned by the loggly endpoint.
// The value for Response should be "ok". Unmarshalling is optimized by ffjson.
type logglyResponse struct {
	Response string `json:"response"`
}
