package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mozz100/screenserve/subprocess"
)

func main() {
	if len(os.Args) <= 2 {
		log.Fatal("Provide a port and a command")
	}
	port := os.Args[1]
	cmd := os.Args[2]
	log.Printf("Command is '%v'", cmd)

	sbpctx := subprocess.GetSubprocess(cmd)

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
		if sbpctx.Parameter != "" || desiredURL == "" {
			sbpctx.Stop()
		}
		if desiredURL != "" {
			sbpctx.StartWith(desiredURL)
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

	homeHandler := func(w http.ResponseWriter, r *http.Request) {
		cssClass := ""
		if sbpctx.Parameter == "" {
			cssClass = "d-none"
		}
		fmt.Fprintf(w, `
		<!doctype html>
		<html lang="en">
		<head>
			<!-- Required meta tags -->
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
			
			<!-- Bootstrap CSS -->
			<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
			
			<title>Wall-e</title>
		</head>
		<body>
			<div class="container mt-3">
				<p class="%v">Currently showing: <code>%v</code></p>
				<form method='POST' action='/launch/'>
					<div class="form-group">
						<input class="form-control" autofocus="autofocus" placeholder="Enter URL" name="url" type="url" value="%v" />
						<small class="form-text text-muted">Include 'http' or 'https'</small>
					</div>
					
					<button class="btn btn-primary" type='submit'>Fling</button>
					&nbsp;or&nbsp;
					<button class="btn btn-light" name='clear' value='1' type='submit'>Clear</button>
				</form>
			</div>
		</body>
		</html>`, cssClass, sbpctx.Parameter, sbpctx.Parameter)
		log.Println("Responded with home page")
	}

	http.HandleFunc("/launch/", launchHandler)
	http.HandleFunc("/", homeHandler)
	log.Printf("Starting listening on port %v\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
