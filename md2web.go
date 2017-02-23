// Package md2web contains the MD2Web trim.Application.
package md2web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	gourl "net/url"
	"path/filepath"
	"strings"

	"github.com/jwowillo/pack"
	"github.com/jwowillo/trim/application"
	"github.com/jwowillo/trim/controller"
	"github.com/jwowillo/trim/request"
	"github.com/jwowillo/trim/response"
	"github.com/jwowillo/trim/url"
	"github.com/russross/blackfriday"
)

// MD2Web is a trim.Applications which turns directories of markdown files and
// folders into a website.
type MD2Web struct {
	*application.Web
}

// New creates a MD2Web excluding the provided files which has the given host.
func New(h, bf string, excs []string) *MD2Web {
	app := &MD2Web{Web: application.NewWeb()}
	app.RemoveAPI()
	app.ClearControllers()
	set := pack.NewHashSet(pack.StringHasher)
	for _, exc := range excs {
		set.Add(exc)
	}
	s := url.NewBuilder(h).
		SetSubdomain(app.Static().Subdomain()).
		SetPath(app.Static().BasePath())
	if err := app.AddController(newClientController(s, bf, set)); err != nil {
		panic(err)
	}
	return app
}

// NewDebug creates an MD2Web that doesn't cache which has the given host.
func NewDebug(h, bf string, excs []string) *MD2Web {
	cf := application.ClientDefault
	cf.CacheDuration = 0
	app := &MD2Web{
		Web: application.NewWebWithConfig(
			cf,
			application.APIDefault,
			application.StaticDefault,
		),
	}
	app.RemoveAPI()
	app.ClearControllers()
	set := pack.NewHashSet(pack.StringHasher)
	for _, exc := range excs {
		set.Add(exc)
	}
	s := url.NewBuilder(h).
		SetSubdomain(app.Static().Subdomain()).
		SetPath(app.Static().BasePath())
	if err := app.AddController(newClientController(s, bf, set)); err != nil {
		panic(err)
	}
	return app
}

// clientController which renders markdown page's based on request paths.
type clientController struct {
	controller.Bare
	staticBuilder *url.Builder
	baseFolder    string
	excludes      pack.Set
}

// newClientController creates a controller with the given template file and
// base folder.
func newClientController(
	sb *url.Builder,
	bf string,
	excs pack.Set,
) *clientController {
	excs.Add("static")
	excs.Add(".git")
	excs.Add(".gitignore")
	return &clientController{
		staticBuilder: sb,
		baseFolder:    bf,
		excludes:      excs,
	}
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
func (c *clientController) Handle(req *request.Request) response.Response {
	fn := req.URL().Path()
	if fn == "/robots.txt" {
		return response.NewStatic("robots.txt")
	}
	path := filepath.Join(c.baseFolder, buildPath(fn))
	path, err := gourl.QueryUnescape(path)
	hl, err := headerLinks(c.baseFolder, path, c.excludes)
	nl, err := navLinks(path, c.excludes)
	bs, err := content(path)
	proto := "http://"
	if req.TLS() != nil {
		proto = "https://"
	}
	static := c.staticBuilder.SetProtocol(proto).Build()
	args := pack.AnyMap{
		"title":       filepath.Base(fn),
		"static":      static,
		"headerLinks": hl,
		"navLinks":    nl,
		"content": strings.Replace(
			string(bs),
			"{{ static }}",
			static.String(),
			-1,
		),
	}
	if err != nil {
		args["headerLinks"] = []linkPair{{Fake: "/", Real: "/"}}
		args["navLinks"] = nil
		args["content"] = fmt.Sprintf("%s couldn't be served.", fn)
		return response.NewTemplateFromString(
			Template,
			args,
			response.CodeNotFound,
		)
	}
	return response.NewTemplateFromString(Template, args, http.StatusOK)
}

// headerLinks are links to files along the provided path except what is in the
// provided set map mapped to their link text.
func headerLinks(bf, path string, excs pack.Set) ([]linkPair, error) {
	ls := []linkPair{linkPair{Real: "/", Fake: "/"}}
	working := ""
	for _, part := range strings.Split(filepath.Dir(path), "/") {
		if part == "." || part == bf {
			continue
		}
		working = filepath.Join(working, part)
		if excs.Contains(working) {
			return nil, fmt.Errorf("%s excluded", working)
		}
		if part == "main.md" {
			break
		}
		if strings.HasSuffix(part, ".md") {
			part = part[:len(part)-len(".md")]
		} else {
			part += "/"
		}
		ls = append(ls, linkPair{Real: "/" + working + "/", Fake: part})
	}
	return ls, nil
}

// navLinks are links to adjacent markdown files and folders to the provided
// path except what is in the excluded provided set mapped to their link text.
//
// Returns an error if the directory of the given path can't be read.
func navLinks(path string, excs pack.Set) ([]linkPair, error) {
	fs, err := ioutil.ReadDir(filepath.Dir(path))
	if err != nil {
		return nil, err
	}
	var ls []linkPair
	for _, f := range fs {
		fn := f.Name()
		if excs.Contains(fn) || excs.Contains(filepath.Base(fn)) {
			continue
		}
		key := f.Name()
		switch mode := f.Mode(); {
		case mode.IsDir():
			key = key + "/"
		case mode.IsRegular():
			if !strings.HasSuffix(fn, ".md") {
				continue
			}
			if fn == "main.md" {
				continue
			}
		}
		if strings.HasSuffix(key, ".md") {
			key = key[:len(key)-len(".md")]
			fn = fn[:len(fn)-len(".md")]
		}
		ls = append(ls, linkPair{Real: key, Fake: fn})
	}
	return ls, nil
}

// content of file at path.
//
// Returns an error if the file isn't a markdown file.
func content(path string) ([]byte, error) {
	if filepath.Ext(path) != ".md" {
		return nil, fmt.Errorf("%s isn't a markdown file", path)
	}
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return blackfriday.MarkdownCommon(bs), nil
}

// buildPath to markdown file represented by given name.
func buildPath(name string) string {
	path := "." + name
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
    <link rel="icon" href="{{ static }}/favicon.png">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
      * {
         font-family: 'Droid Sans', sans-serif;
         color: #2b2b2b;
         word-wrap: break-word;
      }
      img {
      	max-width: 100%;
      }
      #wrapper {
        max-width: 960px;
        margin: 0 auto;
      }
      p {
        line-height: 1.5em;
      }
      pre {
        border: 2px solid #262626;
        padding: 5px;
        background-color: #fff5e6;
        overflow-x: auto;
      }
      code {
        font-family: 'Droid Sans Mono', monospace;;
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
        font-size: 1.2em;
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
      table {
      	width: 100%;
      }
    </style>
  </head>
  <body>
    <div id="wrapper">
      <header>
      	{% for p in headerLinks %}
      	  <a href="{{ p.Real }}">{{ p.Fake }}</a>
      	{% endfor %}
      </header>
      <nav>
        {% for p in navLinks %}
          <a href="{{ p.Real }}">{{ p.Fake }}</a>
        {% endfor %}
      </nav>
      <section>
        {{ content | safe }}
      </section>
    </div>
  <style scoped>
  @import "//fonts.googleapis.com/css?family=Droid+Sans|Droid+Sans+Mono"
  </style>
  </body>
</html>
`

// linkPair is a pair of a real and a fake link.
type linkPair struct {
	Real, Fake string
}
