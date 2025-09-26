Fork of `golang.org/x/term`
===========================
Forked and modified to work a bit better,
at least with [iTerm2](https://iterm2.com)
and (maybe) [Terminal.app](https://support.apple.com/guide/terminal/welcome/mac)].

# Changes
- Blank lines aren't saved to history
- Repeated lines aren't saved to history
- Option+Left/Right works on macOS with iTerm2 and Terminal.app
- Cooked mode, for disabling TTY-like features

# Go terminal/console support

[![Go Reference](https://pkg.go.dev/badge/github.com/magisterquis/goxterm.svg)](https://pkg.go.dev/github.com/magisterquis/goxterm)

This repository provides Go terminal and console support packages.

## Download/Install

The easiest way to install is to run
`go get -u github.com/magisterquis/goxterm`.
