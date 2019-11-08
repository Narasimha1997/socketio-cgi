package main

import "os"

func main() {
	//test
	argv := os.Args
	initServer(argv[1:])
}
