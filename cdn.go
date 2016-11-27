package md2web

import (
	"fmt"

	"github.com/jwowillo/trim"
	"github.com/jwowillo/trim/responses"
)

// NewCDNApplication creates a new trim.Application which serves from the given
// folder.
func NewCDNApplication() *trim.Application {
	application := trim.NewApplication("cdn")
	application.AddController(newCDNController())
	return application
}

// cdnController is a CDN which serves from a particular directory.
type cdnController struct {
	trim.BareController
}

// newCDNController creates a cdnController.
func newCDNController() *cdnController {
	return &cdnController{}
}

// Path matches any path into a variable called 'name'.
func (c *cdnController) Path() string {
	return "/<name />"
}

// Handle returns the static file located at 'name'.
func (c *cdnController) Handle(r *trim.Request) trim.Response {
	path := fmt.Sprintf("static/%s", r.PathArguments()["name"])
	return responses.NewStaticResponse(path)
}
