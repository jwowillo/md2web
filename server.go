package md2web

import (
	"time"

	"github.com/jwowillo/trim"
	"github.com/jwowillo/trim/decorators"
	"github.com/jwowillo/trim/handlers"
)

// NewServer creates a trim.Server running from the given base folder which uses
// the given template file and serves static files from the given folder.
func NewServer(
	domain string,
	port int,
	excludes []string,
) (*trim.Server, error) {
	server := trim.NewServer(domain)
	server.SetHandle404(handlers.NewHTML404Handler())
	server.AddDecoratorFactory(decorators.NewCacheDecoratorFactory(
		time.Hour,
	))
	server.AddDecoratorFactory(decorators.NewAllowDecoratorFactory(
		[]string{"GET"},
	))
	client, err := NewClientApplication(domain, port, excludes)
	if err != nil {
		return nil, err
	}
	server.AddApplication(client)
	server.AddApplication(NewCDNApplication())
	return server, nil
}
