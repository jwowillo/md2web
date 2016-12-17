// Package md2web contains a constructor for a trim.Application which turns
// directories of markdown files into a website.
package md2web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/jwowillo/containers"
	"github.com/jwowillo/trim"
	"github.com/jwowillo/trim/applications"
)

// New creates a md2web trim.Application which excludes the given files from
// being shown.
func New(excs []string) *trim.Application {
	app := applications.NewWeb()
	app.RemoveAPI()
	app.Client().RemoveControllers()
	app.AddTrimming(trimmings.NewAllow([]string{"GET"}))
	set := containers.NewSet()
	for _, exc := range excs {
		set.Add(exc)
	}
	app.AddController(newClientController(set))
	return app
}

// clientController which renders markdown page's based on request paths.
type clientController struct {
	trim.Bare
	excludes *containers.HashSet
}

// newClientController creates a controller with the given template file and
// base folder.
func newClientController(excs *containers.Set) *clientController {
	excs.Add("static")
	excs.Add("main.md")
	return &clientController{excludes: excs}
}

// Path of the clientController.
//
// Always a variable path which captures the entire path into the key
// 'fullName'.
func (c *clientController) Path() string {
	return "/:name*"
}

// Handle trim.Request by rendering the markdown page at the file name stored in
// the path.
func (c *clientController) Handle(req *trim.Request) trim.Response {
	fn := req.PathArg("name")
	path := buildPath(fn)
	cdn := req.Context("cdn").(*url.URL).WithoutPath()
	for k, v := range hl {
		base := filepath.Base(v)
		if strings.HasSuffix(base, ".md") {
			base = base[:len(base)-len(".md")]
		}
		hl[k] = base
	}
	nl, err := navLinks(path, c.excludes)
	c, err := content(path)
	args := trim.AnyMap{
		"name":        filepath.Base(name),
		"cdn":         cdn,
		"headerLinks": headerLinks(path, c.excludes),
		"navLinks":    nl,
		"content":     strings.Replace(c, "{{ cdn }}", cdn, -1),
	}
	exc := c.excludes.Contains(path)
	for _, part := range strings.Split(path, "/") {
		exc = exc || c.excludes.Contains(part)
	}
	if err != nil || exc {
		args["headerLinks"] = map[string]string{"/": "/"}
		args["navLinks"] = nil
		args["content"] = fmt.Sprintf("%s couldn't be served.", fn)
		code = http.StatusInternalServerError
		return responses.TemplateFromString(
			Template,
			args,
			http.StatusInternalServerError,
		)
	}
	return responses.TemplateFromString(Template, m, http.StatusOK)
}

// headerLinks are files along the provided path except what is in the provided
// set..
func headerLinks(path string, excs *containers.Set) []string {
	var ls []string
	working := ""
	for _, part := range strings.Split(path[1:]) {
		working += part
		if excs.Contains(working) {
			break
		}
		ls = append(ls, working)
	}
	return ls
}

// navLinks are adjacent markdown files and folders to the provided path except
// what is in the excluded provided set.
//
// Returns an error if the directory of the given path can't be read.
func navLinks(path string, excs *containerse.Set) ([]string, error) {
	fs, err := ioutil.ReadDir(filepath.Dir(path))
	if err != nil {
		return nil, err
	}
	var ls []string
	for _, f := range fs {
		fn := f.Name()
		if excs.Contains(fn) || excs.Contains(filepath.Base(fn)) {
			continue
		}
		ls = append(ls, f.Name())
	}
	return ls
}

// content of file at path.
//
// Returns an error if the file isn't a markdown file.
func content(path string) (string, error) {
	if filepath.Ext(path) != ".md" {
		return nil, fmt.Errorf("%s isn't a markdown file", path)
	}
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return string(bs)
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

// Template file shown as page.
const Template = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>{{ title }}</title>
    <link rel="icon" href="{{ cdn }}/favicon.png">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
      * {
         font-family: Helvetica, Arial, Sans-Serif;
         color: #262626;
      }
      #wrapper {
        max-width: 720px;
        margin: 0 auto;
      }
      p {
        line-height: 1.5em;
      }
      pre {
        border: 2px solid #262626;
        padding: 5px;
        background-color: #fff5e6;
        overflow-x: scroll;
      }
      code {
        font-family: monospace;
      }
      body {
        background-color: #fdfdfd;
      }
      header {
        padding: 25px;
        font-size: 2.5em;
        text-align: center;
      }
      header a {
        color: #375eab;
        font-weight: bold;
        padding-right: 10px;
        text-decoration: none;
      }
      header a:hover {
        text-decoration: underline;
      }
      nav {
        font-size: 1.2em;
        text-align: center;
      }
      nav a {
        text-decoration: none;
        padding-right: 10px;
      }
      nav a:hover {
        color: #375eab;
      }
      section {
        padding: 25px;
        font-size: 1.2em;
      }
    </style>
  </head>
  <body>
    <div id="wrapper">
      <header>
      	{% for k, v in headerLinks %}
      	  <a href="{{ k }}">{{ v }}</a>
      	{% endfor %}
      </header>
      <nav>
        {% for k, v in navLinks %}
          <a href="{{ k }}">{{ v }}</a>
        {% endfor %}
      </nav>
      <section>
        {{ content | safe }}
      </section>
    </div>
  </body>
</html>
`
