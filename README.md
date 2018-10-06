# Sample site map based web site

A step-by-step guide is here https://golang.org/doc/articles/wiki/

This sample uses a struct to describe each page:
- its base url
- its handler which defaults to the home page handler
- the html template to serve the page. When no file name is provided, it defaults to the page URL.
- CSS style sheet is loaded for the site
- templates are loaded once during init phase (Std mode)

Some features:
- Style sheet is common.
- Inconsistencies are not breaking the site.
- Logs contain all handled exceptions related to the structure.

`Go get` it and start locally using `dev_appserver.py`
The sample is configured to use standard mode of `appengine`.

Known issues:
- Tests are broken because init() fails when starting tests.
- TODO's