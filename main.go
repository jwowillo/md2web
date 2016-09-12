package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/jwowillo/trim"
	"github.com/jwowillo/trim/decorators"
	"github.com/jwowillo/trim/handlers"
	"github.com/jwowillo/trim/responses"
	"github.com/russross/blackfriday"
)

var base string

var templatePath = ".md2web.template.html"

const (
	mainName  = "main"
	nameKey   = "name"
	extension = ".md"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	base = filepath.Base(cwd)
	templatePath = addHome(templatePath)
	server := trim.NewServer()
	server.SetHandle404(handlers.HandleHTML404)
	server.AddDecorator(decorators.CacheDecorator(time.Hour))
	server.AddDecorator(decorators.AllowDecorator([]string{"GET"}))
	server.AddApplication(newApplication())
	server.Run(5000)
}

func newApplication() *trim.Application {
	application := trim.NewApplication("")
	application.AddController(newController())
	application.SetHandle404(handlers.HandleHTML404)
	return application
}

type controller struct {
	trim.BareController
}

func newController() *controller {
	return &controller{}
}

// Path ...
func (c *controller) Path() string {
	return fmt.Sprintf("/<%s />", nameKey)
}

// Handle ...
func (c *controller) Handle(request *trim.Request) trim.Response {
	name := request.PathArguments()[nameKey]
	return renderPage(request, name, "", trim.Code(http.StatusOK))
}

func buildPath(name string) string {
	path := name
	if path == "" || path[len(path)-1] == '/' {
		path += mainName
	}
	path += extension
	return path
}

func renderPage(
	request *trim.Request,
	name, message string,
	code trim.Code,
) trim.Response {
	content, err := ioutil.ReadFile(buildPath(name))
	if err != nil {
		if name == "" {
			return handlers.HandleHTML404(request)
		}
		return renderPage(
			request,
			"",
			fmt.Sprintf("Page at '/%s' doesn't exist.", name),
			trim.Code(http.StatusNotFound),
		)
	}
	result := string(blackfriday.MarkdownCommon(content))
	return responses.NewTemplateResponse(
		templatePath,
		responses.TemplateArguments{
			"base":    base,
			"title":   linkify(name),
			"message": message,
			"links":   links(name),
			"content": result,
		},
		code,
	)
}

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

func links(name string) []string {
	pattern := "<a href=\"%s\">%s</a>"
	path := "./" + name
	path = path[:strings.LastIndex(path, "/")+1]
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil
	}
	var links []string
	width := len(extension)
	for _, file := range files {
		if file.Name() != "main.md" {
			link := "/" + name
			target := file.Name()
			if target[len(file.Name())-width:] == extension {
				target = target[:len(file.Name())-width]
				link += target
			} else {
				link += target + "/"
			}
			links = append(
				links,
				fmt.Sprintf(pattern, link, target),
			)
		}
	}
	return links
}

func addHome(path string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir + "/" + path
}
