package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	subprocess "github.com/mozz100/screenserve/subprocess"
)

// HomeHandler handles GET requests for /
func HomeHandler(sbpctx *subprocess.Context) func(http.ResponseWriter, *http.Request) {
	indexTmpl := template.Must(template.ParseFiles("templates/index.html"))

	homeHandler := func(w http.ResponseWriter, r *http.Request) {
		indexTmpl.Execute(w, sbpctx)
		log.Println("Responded with home page")
	}
	return homeHandler
}

// LaunchHandler handles POST requests
func LaunchHandler(sbpctx *subprocess.Context) func(http.ResponseWriter, *http.Request) {
	launchHandler := func(w http.ResponseWriter, r *http.Request) {
		var desiredURL string
		var slackAPI bool

		r.ParseForm()
		if url, ok := r.PostForm["url"]; ok {
			desiredURL = url[0]
		}
		if url, ok := r.PostForm["text"]; ok {
			desiredURL = url[0]
			slackAPI = true
		}
		if clear, ok := r.PostForm["clear"]; ok && clear[0] == "1" {
			desiredURL = ""
		}

		log.Printf("Launch page with desiredURL = '%v'. Slack: %v\n", desiredURL, slackAPI)
		if desiredURL != "" {
			defer sbpctx.StartWith(desiredURL)
		}
		if sbpctx.Parameter != "" || desiredURL == "" {
			defer sbpctx.Stop()
		}

		if slackAPI {
			log.Println("Responding with plain text 200")
			if desiredURL == "" {
				fmt.Fprint(w, "Cleared!")
			} else {
				fmt.Fprintf(w, "OK, showing %v :tv: :eyes:", desiredURL)
			}

		} else {
			log.Println("Redirecting to home page")
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
	return launchHandler
}
