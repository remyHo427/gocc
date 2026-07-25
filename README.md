# My C compiler
A C compiler written in Go targeting x86-64 assembly,
Currently WIP.

## Dependencies
- Go (1.26.4 tested on linx/amd64)
- GNU Make (4.4.1)
- GCC (16.1.1. This is optional, used as reference in automated test script)

## Supported Commands
- `make`: builds a 
binary named `mycc` under `<PROJECT_ROOT>/bin`.
- `make test`: runs all unit test code, only displays packages with test code in them.
- `make compile_test`: **REQUIRES GCC!**, compile all the C code under `<PROJECT_ROOT>/tests` with both `mycc` and gcc, and then compare the status code returned by each version.