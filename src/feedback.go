package src

/* Running on appengine standard, no main and no go run.
only dev_appserver.py is usable to start the service inside app.yaml directory. */

import (
	"bytes"
	"context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const (
	// webDirectory = "web/" // sitemap>dev_appserver.py app.yaml
	cssFile = "styles.css"
	// TODO load favicon using template
	// Favicon produced using https://realfavicongenerator.net/ (others available)
	faviconFile = "favicon.ico"
)

var webDirectory = "web/" // sitemap>dev_appserver.py app.yaml

// Site map added value:
// - readability
// - flexibility
// - ease of testing which does not need to be tailored to the structure of the site.
var sitemap = []struct {
	url      string
	handler  func(w http.ResponseWriter, r *http.Request)
	filename string // defaults to "" and is searched in webdirectory
	devapp   bool   // defaults to false
}{
	{
		url:     "handler",
		handler: contextlog,
		// filename default is handler name
		devapp: true, // offline test will crash otherwise
	},
	{
		url: "",
		// default handler is root()
		filename: "homepage",
	},
	{
		url:      "sign",
		handler:  sign,
		filename: "feedback",
	},
	{
		url:      "feedback", // synonym of default handler
		filename: "homepage",
	},
}

// A map contains the templates to load them once.
// It can't be included in the sitemap because of looping reference during build.
// The map index is the base url of the page.
var siteTmpl = make(map[string]*template.Template)

// Holding CSS required because of code injection defense
var styleTag template.CSS

func initStylesheet() {
	// Loading style sheet
	tmpl, err := template.ParseFiles(webDirectory + cssFile)
	if err != nil {
		println("init:", err.Error()) // No context and thus no logging
		println("stylesheet is empty")
		return
	}
	// the value returned by ParseGlob.
	var myStyle bytes.Buffer
	if err := template.Must(tmpl, err).Execute(&myStyle, nil); err != nil {
		println("init: ", cssFile, "didn't load", err)
		return
	}
	styleTag = template.CSS(myStyle.String())
}

// In Std mode, loading is during init without context
func initHTMLTemplate(n string) *template.Template {
	pattern := filepath.Join(webDirectory + n + ".html")
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		println(err)
	}
	return template.Must(tmpl, err)
}

// Target is std, i.e. init() and no main.
func init() {
	// os.IsNotExist(err) fails as the error return is "FindFirstFile web/: The parameter is incorrect."
	// Issue is probably related to dev_appserver.py and Python
	if _, err := os.Stat(webDirectory); err != nil {
		// Tests files are in src directory (or a "test" directory)
		webDirectory = "../web/"
	}
	// Registering handlers
	for _, h := range sitemap {
		if h.handler != nil {
			http.HandleFunc("/"+h.url, h.handler)
		} else {
			http.HandleFunc("/"+h.url, root) // default handler is root()
		}
	}
	initStylesheet()
	// Loading templates. It would be ideal to load templates when requested.
	// Due to build issue reported as looping reference, it does not work.
	for _, t := range sitemap {
		if f := t.filename; f == "" {
			// using handler name as filename
			siteTmpl[t.url] = initHTMLTemplate(t.url)
		} else {
			siteTmpl[t.url] = initHTMLTemplate(t.filename)
		}
	}
}

func contextlog(w http.ResponseWriter, r *http.Request) {
	var ctx context.Context
	var module, instance string
	appengine.InstanceID()
	// test crashes as NewContext crashes because it can't find metadata
	ctx = appengine.NewContext(r)
	if ctx != nil {
		module = appengine.ModuleName(ctx)
		instance = appengine.InstanceID()
	}

	data := struct {
		Style      template.CSS
		Module     string
		Instance   string
		Logentries []logentry
	}{
		Style:      styleTag,
		Module:     module,
		Instance:   instance,
		Logentries: displaylog(r),
	}
	var err error
	u := filepath.Base(r.URL.Path)
	// In case of template error, nil is returned
	if t := siteTmpl[u]; t != nil {
		err = t.Execute(w, data)
	} else {
		// request received landed on default page and cannot be served. Logged as info and not error
		pc, _, _, _ := runtime.Caller(0) // Provides the name of the handler
		if filepath.Base(r.URL.Path) != "\\" && appengine.IsDevAppServer() {
			log.Infof(appengine.NewContext(r), "%s was called using %s",
				runtime.FuncForPC(pc).Name(),
				filepath.Base(r.URL.Path))
		}
	}
	if err != nil { // Template execution failed
		log.Errorf(ctx, "%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// You can't loop in root on site map
func root(w http.ResponseWriter, r *http.Request) {
	// No context during init to log
	data := struct {
		Style   template.CSS
		Content string
	}{
		Style:   styleTag,
		Content: r.FormValue("content"),
	}

	var err error
	u := filepath.Base(r.URL.Path)
	if u == "\\" { // Home page base URL is \ is set to "" for consistency
		u = ""
	}
	// A nil template can be:
	// - inconsistent sitemap
	// - invalid template (error during load)
	// - no template provided
	// - appengine internal requests like favicon, styles
	if t := siteTmpl[u]; t != nil {
		// synonym is not reported but only served
		err = t.Execute(w, data)
	} else {
		// request received landed on default page and cannot be served. Logged as info and not error
		pc, _, _, _ := runtime.Caller(0) // Provides the name of the handler
		if filepath.Base(r.URL.Path) != "\\" && appengine.IsDevAppServer() {
			log.Infof(appengine.NewContext(r), "%s was called using %s",
				runtime.FuncForPC(pc).Name(),
				filepath.Base(r.URL.Path))
		}
	}
	if err != nil {
		log.Errorf(appengine.NewContext(r), "%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func sign(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Style   template.CSS
		Content string
	}{
		Style:   styleTag,
		Content: r.FormValue("content"),
	}

	var err error
	if t := siteTmpl[filepath.Base(r.URL.Path)]; t != nil {
		err = t.Execute(w, data)
	}
	if err != nil {
		log.Errorf(appengine.NewContext(r), "%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
