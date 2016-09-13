package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/jwowillo/trim"
	"github.com/jwowillo/trim/decorators"
	"github.com/jwowillo/trim/handlers"
	"github.com/jwowillo/trim/responses"
	"github.com/russross/blackfriday"
)

// newServer creates a trim.Server running from the given base folder which uses
// the given template file.
func newServer(base, template string) *trim.Server {
	server := trim.NewServer()
	server.SetHandle404(handlers.HandleHTML404)
	server.AddDecorator(decorators.CacheDecorator(time.Hour))
	server.AddDecorator(decorators.AllowDecorator([]string{"GET"}))
	server.AddApplication(newApplication(base, template))
	return server
}

// newApplication creates a md2web trim.Application running from the given base
// folder which uses the given template file.
func newApplication(base, template string) *trim.Application {
	application := trim.NewApplication("")
	application.AddController(newController(base, template))
	return application
}

// controller which renders markdown page's based on request paths.
type controller struct {
	trim.BareController
	base, template string
}

// newController creates a controller with the given template file and base
// folder.
func newController(base, template string) *controller {
	return &controller{template: template, base: base}
}

// Path of the controller.
//
// Always a variable path which captures the entire path into the key 'name'.
func (c *controller) Path() string {
	return "/<name />"
}

// Handle trim.Request by rendering the markdown page at the file name stored in
// the path.
func (c *controller) Handle(request *trim.Request) trim.Response {
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
func (c *controller) renderPage(
	name, message string,
	code trim.Code,
) (trim.Response, error) {
	if filepath.Base(name) == "main" {
		return c.renderPage(
			"",
			fmt.Sprintf("Page at '/%s doesn't exist.", name),
			trim.Code(http.StatusNotFound),
		)
	}
	content, err := ioutil.ReadFile(buildPath(name))
	// If request is for a main page, re render with error message.
	if err != nil {
		if name == "" {
			return nil, errors.New("Can't read file.")
		}
		return c.renderPage(
			"",
			fmt.Sprintf("Page at '/%s' doesn't exist.", name),
			trim.Code(http.StatusNotFound),
		)
	}
	return responses.NewTemplateResponse(
		c.template,
		responses.TemplateArguments{
			"base":    c.base,
			"title":   linkify(name),
			"message": message,
			"links":   links(name),
			"content": string(blackfriday.MarkdownCommon(content)),
		},
		code,
	), nil
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

// links returns links to all of the sibling pages to the page at the given
// name.
func links(name string) []string {
	pattern := "<a href=\"%s\">%s</a>"
	path := "./" + name
	path = path[:strings.LastIndex(path, "/")+1]
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil
	}
	name = name[:strings.LastIndex(name, "/")+1]
	var links []string
	links = append(links, fmt.Sprintf(pattern, path[1:], "."))
	for _, file := range files {
		if file.Name() != "main.md" {
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

// buildPath to markdown file represented by given name.
func buildPath(name string) string {
	path := name
	if path == "" || path[len(path)-1] == '/' {
		path += "main"
	}
	path += ".md"
	return path
}
