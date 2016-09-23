package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/jwowillo/trim"
	"github.com/jwowillo/trim/handlers"
	"github.com/jwowillo/trim/responses"
	"github.com/russross/blackfriday"
)

// newClientApplication creates a md2web trim.Application running from the given
// base folder which uses the given template file.
func newClientApplication(
	base, template, static string,
	port int,
) *trim.Application {
	application := trim.NewApplication("")
	application.AddController(newClientController(
		base,
		template,
		static,
		port,
	))
	return application
}

// clientController which renders markdown page's based on request paths.
type clientController struct {
	trim.BareController
	base, template, static string
	port                   int
}

// newClientController creates a controller with the given template file and
// base folder.
func newClientController(
	base, template, static string,
	port int,
) *clientController {
	return &clientController{
		template: template,
		base:     base,
		static:   static,
		port:     port,
	}
}

// Path of the clientController.
//
// Always a variable path which captures the entire path into the key 'name'.
func (c *clientController) Path() string {
	return "/<name />"
}

// Handle trim.Request by rendering the markdown page at the file name stored in
// the path.
func (c *clientController) Handle(request *trim.Request) trim.Response {
	name := request.PathArguments()["name"]
	response, err := c.renderPage(name, "", trim.Code(http.StatusOK))
	if err != nil {
		return handlers.HandleHTML404(request)
	}
	return response
}

// renderPage based on the name of the markdown file, the message to display on
// the page, and code meant for the controller to return.
//
// The trim.Request is passed to handle
func (c *clientController) renderPage(
	name, message string,
	code trim.Code,
) (trim.Response, error) {
	if filepath.Base(name) == c.static {
		return c.errorResponse(name)
	}
	if filepath.Base(name) == "main" {
		return c.errorResponse(name)
	}
	content, err := ioutil.ReadFile(buildPath(name))
	if err != nil {
		if name == "" {
			return nil, errors.New("Can't read file.")
		}
		return c.errorResponse(name)
	}
	return responses.NewTemplateResponse(
		c.template,
		responses.TemplateArguments{
			"base":    c.base,
			"title":   linkify(name),
			"message": message,
			"links":   c.links(name),
			"port":    c.port,
			"content": string(blackfriday.MarkdownCommon(content)),
		},
		code,
	), nil
}

// errorResponse returns the result of calling renderPage with not-found data.
func (c *clientController) errorResponse(name string) (trim.Response, error) {
	return c.renderPage(
		"",
		fmt.Sprintf("Page at '/%s doesn't exist.", name),
		trim.Code(http.StatusNotFound),
	)
}

// links returns links to all of the sibling pages to the page at the given
// name.
func (c *clientController) links(name string) []string {
	pattern := "<a href=\"%s\">%s</a>"
	path := "./" + name
	path = path[:strings.LastIndex(path, "/")+1]
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil
	}
	name = name[:strings.LastIndex(name, "/")+1]
	var links []string
	for _, file := range files {
		if file.Name() != "main.md" && file.Name() != c.static {
			link := "/" + name
			target := file.Name()
			if target[len(file.Name())-3:] == ".md" {
				target = target[:len(file.Name())-3]
				link += target
				links = append(
					links,
					fmt.Sprintf(pattern, link, target),
				)
			} else if file.IsDir() {
				link += target + "/"
				links = append(
					links,
					fmt.Sprintf(pattern, link, target),
				)
			}
		}
	}
	return links
}

// linkify convertes a name into a name where each subpath is linked to its
// page.
func linkify(name string) string {
	pattern := "<a href=\"%s\">%s</a>"
	name = "/" + name
	link := ""
	i := 0
	var j int
	for j = 0; j < len(name); j++ {
		if name[j] == '/' {
			link += fmt.Sprintf(pattern, name[:j+1], name[i:j+1])
			i = j + 1
		}
	}
	link += fmt.Sprintf(pattern, name[:j], name[i:j])
	return link
}

// buildPath to markdown file represented by given name.
func buildPath(name string) string {
	path := name
	if path == "" || path[len(path)-1] == '/' {
		path += "main"
	}
	path += ".md"
	return path
}
