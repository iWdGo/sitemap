# Sample site map based web site

A step-by-step guide is here https://golang.org/doc/articles/wiki/

This sample uses a struct to describe each page:
- its base url
- its handler which defaults to the home page handler
- the html template to serve the page. When no file name is provided, it defaults to the page URL.
- CSS style sheet is loaded for the site
- templates are loaded once during init phase (Std mode)
- static files are only a filename with extension

Some features:
- Style sheet is common.
- Inconsistencies are not breaking the site.
- Static files folder cannot be browsed.
- Logs contain all handled exceptions related to the structure.
- Basic tests

`Go get` it and start locally using `dev_appserver.py` or `dev_appserver.py %CD%`.
Testing is executed using `go test ./src`.
This sample is configured to use standard mode of `appengine`.

Known issues:
- favicon is not loaded
- To deploy on Google Cloud, app.yaml must be moved inside the `src` directory.