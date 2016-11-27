// Package main runs a server which serves markdown files and folders as a
// website.
package main

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/jwowillo/md2web"
)

// main runs the server on the given port or 5000 by default.
func main() {
	domain, err := getDomain()
	port, err := getPort()
	if err != nil {
		log.Fatal(err)
	}
	server, err := md2web.NewServer(domain, port, []string{"README.md"})
	if err != nil {
		log.Fatal(err)
	}
	server.Run(port)
}

// getDomain to listen at.
func getDomain() (string, error) {
	if len(os.Args) != 3 {
		return "", errors.New("Usage: md2web <domain> <port>")
	}
	return os.Args[1], nil
}

// getPort to listen at.
func getPort() (int, error) {
	if len(os.Args) != 3 {
		return -1, errors.New("Usage: md2web <domain> <port>")
	}
	portArg := os.Args[2]
	portVal, err := strconv.Atoi(portArg)
	if err != nil {
		return -1, err
	}
	return portVal, nil
}
