# Gelf Handler
Adds the [GELF](http://docs.graylog.org/en/2.1/pages/gelf.html) format for Graylog-based logging.
GELF can be udp+tcp based, and supports chunking with udp, thus avoiding reconnection- and performance issues.

# Duplicate keys
Currently log15 will duplicat keys to the contect list. Gelf expects a map, therefore keys have to be unique.
This implementation assures that the last value is used for this key.

```go
    l1:=log.New("foo","bar")
    l1.Info("a message","foo","baz")
    // Output: in GELF: msg: "a message", foo: "baz"
```

# Limitations
- only supports udp with gzip compression.
- buffer size not adjustable