package src

import (
	"io/ioutil"
	"net/http"
	"testing"
)

var url = "http://localhost:8080"

// Tests online using dev_appserver.py
func TestHandlersOnline(t *testing.T) {
	// Test fail if site is not locally online. If deployed in cloud (gcloud), address is invalid
	client := &http.Client{}

	for _, h := range sitemap {

		req, err := http.NewRequest("GET", url+"/"+h.page, http.NoBody)
		if err != nil {
			t.Fatal(err)
		}

		r, err := client.Do(req)
		if err != nil {
			t.Fatal("client do:", err)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal("page received failed with ", err)
		}
		if len(body) == 0 {
			t.Fatal("empty page returned")
		}
		req.Body.Close()
	}
}
