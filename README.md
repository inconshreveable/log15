# log15

Package log15 provides an opinionated, simple toolkit for best-practice logging that is both human and machine readable. It is modeled after the standard library's io and net/http packages.

## Features
- A simple, easy-to-understand API
- Promotoes structured logging by encouraging use of key/value pairs
- Child loggers which inherit and add their own private context
- Lazy evaluation of expensive operations
- Simple Handler interface allowing for construction of flexible, custom logging configurations with a tiny API.
- Color terminal support
- Built-in support for logging to files, streams, syslog, and the network
- Support for forking records to multiple handlers, buffering records for output, failing over from failed handler writes, + more

## Documentation

The package documentation is extensive and complete. Browse on godoc:

#### [log15 API Documentation](https://godoc.org/github.com/inconshreveable/log15)

## License
Apache
