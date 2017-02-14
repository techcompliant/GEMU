[![GoDoc](https://godoc.org/gemu.techcompliant.com/gemu?status.svg)](https://godoc.org/gemu.techcompliant.com/gemu)

# Introduction

GEMU is the DCPU emulator that powers Tech Compliant.
In order to allow DCPU developers to write software targeting our DCPU systems easier, we are open sourcing the core of our emulator, as well as a few small utilities to allow it to be used directly, without requiring developers to build their own tools around our emulator.

# Documentation

GEMU documentation is currently being built, though a very minimal amount of information can be gathered by the simple godoc tool.  This can be accessed via the GoDoc tag above.

# GEMUSingle

Included in this repo is a simple single DCPU emulator.  If you have installed Go correctly, and set up a proper gopath, this can be compiled via `make` either from this main directory, or from in the GEMUSingle directory.  Of course, if you are more comfortable with the `go` tool, feel free to use it directly.

# GEMU Compatible projects

The following is a short list of projects that are confirmed to be working with GEMU.  Note that TC has changed a few device IDs, specifically the LEM and keyboard IDs, so stock DCPU code may not run directly on this emulator.

* [BBOS](https://github.com/madmockers/BareBonesOS)
* [Admiral](https://github.com/orlof/dcpu-admiral)
* [MoonPatrol](https://github.com/orlof/dcpu-moonpatrol)
* [DCPU-MUD](https://github.com/orlof/dcpu-mud)
* [DC-DOS](https://github.com/interfect/bbfs)
* 