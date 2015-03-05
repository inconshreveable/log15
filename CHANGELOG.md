# log15 changelog
-------

This file keeps track of the log15 changes and improvements since `v2.9`.
It also contains information on how to migrate code to new versions when breaking changes are being made.


<!---

Please follow this format:


## vX.Y yyyy-mm-dd

#### Breaking changes

**Renamed `Warn` to `Warning`**

We thought it would be funny to break everyones code, so we changed `Warn(msg string, ctx ...interface{})` to `Warning(ctx interface{}, msg ...string)`.
As you can see, we also re-ordered the arguments. Updating is very simple. Change `logger.Warn("some message", yourCtx)` to `logger.Warning(yourCtx, "some message")

#### Major changes & new features

**Added something**

And the explanation goes here

**Changed something**

And some more explanation

#### Minor changes
 - list of
 - small changes
 - that were made

-->

## upcoming (master)

#### Major changes & new features

**Add support for Loggly**

`LogglyHandler` was added in the subpackage `ext/loggly`. This handler sends logs to the [Loggly](http://loggly.com/) logging service.
For more information about the handler read the godoc: https://godoc.org/gopkg.in/inconshreveable/log15.v2/ext/loggly

#### Minor changes
 - none yet
