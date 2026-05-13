# Changelog

## v3.2.1 - 2026-05-13

- Restored `v3/ext.RandId`, which was accidentally omitted from v3.1.0 and
  could break upgrades from `v3.0.0-testing.5`. Thanks to Jonathan Stacks for
  the fix.

## v3.2.0 - 2026-05-13

- Removed the `github.com/mattn/go-colorable` dependency from the root and
  `v3` modules. This also removes the transitive
  `github.com/mattn/go-isatty` dependency.
- Switched terminal detection to `golang.org/x/term.IsTerminal`.
- Kept the deprecated root `term` package available for source compatibility.
  `term.IsTty` now delegates to `x/term`, and the exported `Termios` types are
  still present but marked deprecated.
- Updated `compat` to require the root module version that no longer depends on
  `go-colorable`, allowing `go mod tidy` to remove `go-colorable`,
  `go-isatty`, and stale `x/sys` checksum entries from `compat`.
- Updated module metadata and CI configuration to current Go and `x/sys`
  versions.

## v3.1.0 - 2025-07-17

- Prepared the `github.com/inconshreveable/log15/v3` module for release.
- Added the `compat` module for adapting between legacy log15 interfaces and
  the v3 APIs.
- Improved compatibility coverage for handlers, loggers, record conversion, and
  key name handling.
- Preserved the stack-backed caller record work from Nikolay Petrov, including
  allocation and benchmark follow-ups.

## v3.0.0-testing.5 - 2022-11-29

- Added the initial `github.com/inconshreveable/log15/v3` module layout.
- Added compatibility wrappers for legacy handlers and loggers.
- Updated import paths and formatting for the v3 module.
- Thanks to Josh Robson Chase for the v3 compatibility work.

## v2.16.0 - 2022-11-21

- Switched root terminal detection to `golang.org/x/term`.
- Added GitHub Actions CI.
- Added CI coverage for `ppc64le`, from Devendranath Thadi.

## v2.15-no-caller - 2020-08-03

- Added a variant that removes caller reporting.
- Thanks to Nikolay Petrov for the caller and stack record work that led to this
  release path.

## v2.15 - 2020-01-09

- Made `LvlFromString` more flexible.
- Updated CI coverage for newer Go versions.
