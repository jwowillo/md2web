// Package main runs a server containing the md2web trim.Application.
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jwowillo/md2web"
	"github.com/jwowillo/trim/server"
)

var (
	url  string
	port int
)

// main runs the server on the given url and port.
func main() {
	h := url
	if url == "localhost" || port != 80 {
		h += fmt.Sprintf(":%d", port)
	}
	server.New(url, port).Serve(md2web.New(h, "", []string{"README.md"}))
}

// init parses the URL and port from the command line.
func init() {
	message := []byte("Usage: md2web <url> <port:int>\n")
	if len(os.Args) != 3 {
		log.Fatal(message)
	}
	host = os.Args[1]
	portArg := os.Args[2]
	portVal, err := strconv.Atoi(portArg)
	if err != nil {
		log.Fatal(err)
	}
	port = portVal
}
