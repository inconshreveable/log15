// provides a Gelf Handler for log15
// taken from https://gist.github.com/jmtuley/d4b09617967e59c58c3e
// and parts from https://github.com/gemnasium/logrus-hooks
// it has no external dependencies outside the std library and of course log15
// the actual Handler still resides inside the log15 package; here we have the supporting functions
package gelf
