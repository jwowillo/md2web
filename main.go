// Package main runs a server which serves markdown files and folders as a
// website.
package main

import (
	"errors"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

// TODO: Put applications in own packages so they are black-box testable.

const (
	// template file name.
	template = ".md2web.template.html"
	// static file folder.
	static = "static"
)

// main runs the server on the given port or 5000 by default.
func main() {
	base, err := getBase()
	if err != nil {
		log.Fatal(err)
	}
	home, err := getHome()
	if err != nil {
		log.Fatal(err)
	}
	domain, err := getDomain()
	port, err := getPort()
	if err != nil {
		log.Fatal(err)
	}
	newServer(domain, base, home+"/"+template, static, port).Run(port)
}

// getBase gets the base folder the trim.Server is running from.
func getBase() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Base(cwd), nil
}

// getHome returns the user's home directory or an error if it can't be
// retrieved.
func getHome() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
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
