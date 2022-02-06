# goshell
My first [Go]("https://go.dev") project! It's built to be compatible with Mac, Linux and Windows (although I only tested it on Linux). The only requirements are a `gzip` command that can be used from your shell. (Of course you will also need [Go]("https://go.dev"))

# Usage
You can use it to test some [Go]("https://go.dev") code without all the hastle of creating a file and running it. After cloning, you can build the binary by running `cd goshell` then running `go build goshell.go`. This will generate a binary that you can run. Provide the `--imports` option to declare packages that you want imported.

# Example
An example of using goshell.
```
$ ./goshell --imports fmt
GoShell v0.0.1
Type '!!help' for more info
>>> fmt.Println("Hello, World!")
>>> !!run
Hello, World!
>>> !!cache
>>> !!exit
$ ./goshell
>>> !!restore
>>> !!run
Hello, World!
```
