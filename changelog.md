Changelog
=========

Branch `betterctrlc`
--------------------
- Added changelog
- Added `ReadWriter`, which combines an `io.Reader` and an `io.Writer` into an
  `io.ReadWriter`.
- Added `StdioRW`, a `ReadWriter` wrapping stdin and stdout.
- Added `CtrlC` and `ErrCtrlC`, returned from `Terminal.ReadLine` on Ctrl+C
  instead of EOF.
- Updated dependencies.
