# gobuildall
Small utility that compiles all packages and all tests in a Go project without keeping the resulting binaries. Useful when added as an external command to the Golang IDE or VS Code if you want a button or a keyboard shortcut that compiles (and which thereby syntax checks) all relevant files in your project (including the tests, that a normal build would have skipped). 

When invoked like this:

`$ gobuildall`

It will build every package and test in your current working directory.

Install it with `go get github.com/janflyborg/gobuildall` and add it as a shortcut to your IDE of choice.
