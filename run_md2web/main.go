// Package main runs a server containing the md2web trim.Application.
package main

import (
	"log"
	"os"
	"strconv"

	"github.com/jwowillo/md2web"
)

var (
	domain string
	port   int
)

// main runs the server on the given domain and port.
func main() {
	server, err := md2web.NewServer(domain, port, []string{"README.md"})
	if err != nil {
		log.Fatal(err)
	}
	server.Run(port)
}

// init parses the domain and port from the command line.
func init() {
	message := []byte("Usage: md2web <domain> <port:int>\n")
	if len(os.Args) != 3 {
		os.Stderr.Write(message)
		os.Exit(1)
	}
	domain = os.Args[1]
	portArg := os.Args[2]
	portVal, err := strconv.Atoi(portArg)
	if err != nil {
		os.Stderr.Write(message)
		os.Exit(1)
	}
	port = portVal
}
