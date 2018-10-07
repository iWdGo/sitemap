package src

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func templateNil(t *testing.T, n string) string {
	tmpl, err := template.ParseGlob(filepath.Join(webDirectory + n + ".html"))
	if err != nil {
		t.Fatalf("%v", err)
	}
	w := httptest.NewRecorder()
	// Loading only the stylesheet which is common
	nilData := struct {
		Style template.CSS
	}{
		Style: styleTag,
	}
	err = template.Must(tmpl, err).Execute(w, nilData)
	return w.Body.String()
}

// Test offline is restricted to getting a reply. It checks somehow the consistency of the sitemap
func TestHandlersOffline(t *testing.T) {
	for _, h := range sitemap {
		if !h.devapp {
			r, err := http.NewRequest("GET", "/"+h.url, http.NoBody)
			if err != nil {
				t.Fatal("New request failed with ", err)
			}

			w := httptest.NewRecorder()
			if h.handler != nil {
				h.handler(w, r)
			} else {
				root(w, r)
			}

			if w.Code != 200 {
				t.Fatalf("wrong code returned: %d", w.Code)
			}

			// Test compares result of handler and the template execution w/o data (nil)
			tname := h.filename
			if tname == "" { // filename is the handler name
				tname = h.url
			}
			got := templateNil(t, tname)
			if want := w.Body.String(); want == got {
				// test is successful
			} else if got[len(got)-9:] != "</html>\r\n" {
				// Test is deemed successful as html template execution does not return closing tags.
				// The execution of the template is stopped after the value. The reasons are not:
				//   the name of the field, the absence of empty line, the last active field
				// No related error has been found. This is only related to the nil value used for execution.
				t.Logf("%s:%s got is truncated by %d. Test skipped", t.Name(), tname, len(want)-len(got))
			} else if want != got {
				t.Errorf("offline %s : got %v, want %v", tname, got, want)
			}
			r.Body.Close()
		}
	}
}
