# run
This package wraps the interpreter and runs it against given source file for which file path is given as input argument.

## Usage
1. run `make build`
2. add symlink from `.../src/run/flow`* to `/usr/local/bin`
3. flow is now usable from the terminal
4. run `flow {{relative_filename}}`

*enter full path to the flow executable here

When this works, to update the interpreter only run `make build` again