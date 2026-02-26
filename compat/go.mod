module github.com/inconshreveable/log15/compat

go 1.24.0

require (
	// can't use "replace" directive for this import which is just the parent module
	github.com/inconshreveable/log15 v2.16.0+incompatible
	github.com/inconshreveable/log15/v3 v3.0.0-testing.5
)

require (
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/term v0.40.0 // indirect
)

replace github.com/inconshreveable/log15/v3 => ../v3
