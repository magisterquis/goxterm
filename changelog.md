Changelog
=========

`v0.0.1-beta.3`
---------------
- Added changelog
- Added cooked mode.
- Added `ReadWriter`, which combines an `io.Reader` and an `io.Writer` into an
  `io.ReadWriter`.
- Added `StdioRW`, a `ReadWriter` wrapping stdin and stdout.
- Added `CtrlC` and `ErrCtrlC`, returned from `Terminal.ReadLine` on Ctrl+C
  instead of EOF.
- Added `CookedEscapeCodes` to get empty escape codes.
- Updated dependencies and merged the latest upstream.
