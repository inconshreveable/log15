# Gelf Handler
Adds the [GELF](http://docs.graylog.org/en/2.1/pages/gelf.html) format for Graylog-based logging.
GELF can be udp+tcp based, and supports chunking with udp, thus avoiding reconnection- and performance issues.

# Limitations
- only supports udp with gzip compression.
- buffer size not adjustable