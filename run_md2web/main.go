// Package main runs a server containing the md2web trim.Application.
package main

import (
	"log"
	"os"
	"strconv"

	"github.com/jwowillo/md2web"
	"github.com/jwowillo/trim"
)

var (
	domain string
	port   int
)

// main runs the server on the given domain and port.
func main() {
	trim.NewServer(domain, port).Serve(md2web.New([]string{"README.md"}))
}

// init parses the domain and port from the command line.
func init() {
	message := []byte("Usage: md2web <domain> <port:int>\n")
	if len(os.Args) != 3 {
		log.Fatal(message)
	}
	domain = os.Args[1]
	portArg := os.Args[2]
	portVal, err := strconv.Atoi(portArg)
	if err != nil {
		log.Fatal(err)
	}
	port = portVal
}
