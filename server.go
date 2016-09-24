package main

import (
	"time"

	"github.com/jwowillo/trim"
	"github.com/jwowillo/trim/decorators"
	"github.com/jwowillo/trim/handlers"
)

// newServer creates a trim.Server running from the given base folder which uses
// the given template file and serves static files from the given folder.
func newServer(domain, base, template, static string, port int) *trim.Server {
	server := trim.NewServer(domain)
	server.SetHandle404(handlers.HandleHTML404)
	server.AddDecorator(decorators.CacheDecorator(time.Hour))
	server.AddDecorator(decorators.AllowDecorator([]string{"GET"}))
	server.AddApplication(newClientApplication(
		base,
		template,
		static,
		domain,
		port,
	))
	server.AddApplication(newCDNApplication(static))
	return server
}
