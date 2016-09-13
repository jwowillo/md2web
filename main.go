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

const (
	// template file name.
	template = ".md2web.template.html"
	// defaultPort for server to listen at.
	defaultPort = 5000
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
	port, err := getPort(defaultPort)
	if err != nil {
		log.Fatal(err)
	}
	newServer(base, home+"/"+template, static, port).Run(port)
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

// getPort returns the port to run the server at or the default port if no
// arguments are given to the program.
func getPort(defaultPort int) (int, error) {
	port := defaultPort
	if len(os.Args) > 2 {
		return -1, errors.New("Too many arguments.")
	}
	if len(os.Args) == 2 {
		portArg := os.Args[1]
		portVal, err := strconv.Atoi(portArg)
		if err != nil {
			return -1, err
		}
		port = portVal
	}
	return port, nil
}
