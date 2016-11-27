# md2web
md2web converts a folder hierarchy of markdown files into a navigable website.

The functionality is provided in several forms:
* Runnable: The command `run_md2web` starts an md2web server in the working
  directory.
* Server: Server can be run from other programs when desired.
* Applications: Applications can be attached to other trim servers to provide
  the same functionality in a modular way.

## Documentation:

Documentation is located at https://godoc.org/github.com/jwowillo/md2web.

## Install
Change directory to `run_md2web` and type `make`. The template home page will be
copied into the user's home directory.

## Import

Import with `import "github.com/jwowillo/md2web"`. Make sure the previous
installation instructions were run. Then run `go get`.

## Run

The command `run_md2web <domain> <port:int>` can be used to start the server.
