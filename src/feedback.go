package src

/* Running on appengine standard, no main and no go run.
only dev_appserver.py is usable to start inside app.yaml directory. */

import (
	"bytes"
	"context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"html/template"
	"net/http"
	"path/filepath"
	"runtime"
)

const (
	webDirectory = "../web/" // src>dev_appserver.py
	cssFile      = "styles.css"
	// TODO load favicon using template
	// Favicon produced using https://realfavicongenerator.net/ (others available)
	faviconFile = "favicon.ico"
)

// Site map value is readability and its availability for testing.
// TODO add template in struct to avoid re-loading templates
var sitemap = []struct {
	page     string
	handler  func(w http.ResponseWriter, r *http.Request)
	filename string // defaults to ""
	devapp   bool   // defaults to false
}{
	{
		page:    "handler",
		handler: contextlog,
		devapp:  true, // offline test will crash otherwise
	},
	{
		page:     "",
		handler:  root,
		filename: "homepage",
	},
	{
		page:     "sign",
		handler:  sign,
		filename: "feedback",
	},
	{
		page:     "feedback", // synonym
		handler:  root,
		filename: "homepage",
	},
}

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

// Target is std, i.e. init() and no main
func init() {
	// Registering handlers
	for _, h := range sitemap {
		http.HandleFunc("/"+h.page, h.handler) // displays contextual only if deployed
	}
	initStylesheet()
}

// Loading templates
func tmplLoad(name string, r *http.Request) *template.Template {
	ctx := appengine.NewContext(r)
	//
	// pattern is the glob pattern used to find all the html files.
	pattern := filepath.Join(webDirectory + name + ".html")
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		log.Errorf(ctx, "%v", err)
	}
	// the value returned by ParseGlob.
	return template.Must(tmpl, err)
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
	contextlog := tmplLoad("contextlog", r)
	err := contextlog.Execute(w, data)
	if err != nil {
		log.Errorf(appengine.NewContext(r), "%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// You can't loop in root on site map
func root(w http.ResponseWriter, r *http.Request) {
	homepage := tmplLoad(webDirectory+"homepage", r)

	data := struct {
		Style   template.CSS
		Content string
	}{
		Style:   styleTag,
		Content: r.FormValue("content"),
	}

	err := homepage.Execute(w, data)
	if err != nil {
		log.Errorf(appengine.NewContext(r), "%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	pc, _, _, _ := runtime.Caller(0) // Provides the name of the handler

	if filepath.Base(r.URL.Path) != "\\" && appengine.IsDevAppServer() {
		log.Infof(appengine.NewContext(r), "%s was called using %s",
			runtime.FuncForPC(pc).Name(),
			filepath.Base(r.URL.Path))
	}
}

func sign(w http.ResponseWriter, r *http.Request) {
	feedback := tmplLoad(webDirectory+"feedback", r)
	data := struct {
		Style   template.CSS
		Content string
	}{
		Style:   styleTag,
		Content: r.FormValue("content"),
	}
	err := feedback.Execute(w, data)
	if err != nil {
		log.Errorf(appengine.NewContext(r), "%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
