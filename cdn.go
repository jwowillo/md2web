package main

import (
	"github.com/jwowillo/trim"
	"github.com/jwowillo/trim/responses"
)

// newCDNApplication creates a new trim.Application which serves from the given
// folder.
func newCDNApplication(static string) *trim.Application {
	application := trim.NewApplication("cdn")
	application.AddController(newCDNController(static))
	return application
}

// cdnController is a CDN which serves from a particular directory.
type cdnController struct {
	trim.BareController
	static string
}

// newCDNController creates a cdnController which serves from a given directory.
func newCDNController(static string) *cdnController {
	return &cdnController{static: static}
}

// Path matches any path into a variable called 'name'.
func (c *cdnController) Path() string {
	return "/<name />"
}

// Handle returns the static file located at 'name'.
func (c *cdnController) Handle(r *trim.Request) trim.Response {
	path := c.static + "/" + r.PathArguments()["name"]
	return responses.NewStaticResponse(path)
}
