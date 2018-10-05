# Sample site map based web site

A step-by-step guide is here https://golang.org/doc/articles/wiki/

This sample adds:
- a site map struct to describe each url of the site
- how a handler can dynamically handle synonyms.
- basic tests of handlers is one go test using the map.
- CSS style sheet loaded using template to comply with code injection.
- online log display

web folder contains HTML and the CSS stylesheet.

`Go get` it and start locally using `dev_appserver.py`
The sample is configured to use standard mode of `appengine`.

Some TODO's are included.