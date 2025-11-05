package main

import (
	"flag"
)

func main() {
	port := flag.Int("p", 8000, "port on which the server will listen")
	flag.Parse()

	RunServer(*port)
}
