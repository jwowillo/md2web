// Package md2web contains trim.Applications which turn folder hierarchies of
// markdown files into websites.
package md2web

import (
	"time"

	"github.com/jwowillo/trim"
	"github.com/jwowillo/trim/decorators"
	"github.com/jwowillo/trim/handlers"
)

// NewServer creates a trim.Server with the md2web trim.Applications.
//
// The domain and port of the server are passed along with the files to exclude
// from being shown on the website.
func NewServer(
	domain string,
	port int,
	excludes []string,
) (*trim.Server, error) {
	server := trim.NewServer(domain)
	server.SetHandle404(handlers.NewHTML404Handler())
	for _, f := range []trim.DecoratorFactory{
		decorators.NewCacheDecoratorFactory(time.Hour),
		decorators.NewAllowDecoratorFactory([]string{"GET"}),
	} {
		server.AddDecoratorFactory(f)
	}
	client, err := NewClientApplication(domain, port, excludes)
	if err != nil {
		return nil, err
	}
	server.AddApplication(client)
	server.AddApplication(NewCDNApplication())
	return server, nil
}
