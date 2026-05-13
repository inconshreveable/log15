module github.com/inconshreveable/log15/compat

go 1.24.0

require (
	// can't use "replace" directive for this import which is just the parent module
	github.com/inconshreveable/log15 v0.0.0-20260226202544-23f08a1ae593
	github.com/inconshreveable/log15/v3 v3.1.0
)

require (
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/term v0.40.0 // indirect
)

replace github.com/inconshreveable/log15/v3 => ../v3
