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

	sbpctx := subprocess.Context{}
	ch := *subprocess.Subprocess(cmd, &sbpctx)

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

		log.Printf("Launch page with desiredURL = '%v'. Slack: %v\n", desiredURL, slackAPI)
		if sbpctx.SubProc != nil || desiredURL == "" {
			ch <- subprocess.Instruction{
				Instruction: "stop",
				Parameter:   "",
			}
		}
		if desiredURL != "" {
			ch <- subprocess.Instruction{
				Instruction: "start",
				Parameter:   desiredURL,
			}
		}

		if slackAPI {
			log.Println("Responding with plain text 200")
			fmt.Fprint(w, "OK :tv: :eyes:")
		} else {
			log.Println("Redirecting to home page")
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}

	homeHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
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
			<div class="container">
				<p>Put in a URL to fling it to Wall-e</p>
				<form method='POST' action='/launch/'>
					<div class="form-group">
						<input class="form-control" autofocus="autofocus" name="url" type="url" />
						<small class="form-text text-muted">Include 'http' or 'https'</small>
					</div>
					
					<button class="btn btn-primary" type='submit'>Go</button>
				</form>
			</div>
		</body>
		</html>`)
		log.Println("Responded with home page")
	}

	http.HandleFunc("/launch/", launchHandler)
	http.HandleFunc("/", homeHandler)
	log.Printf("Starting listening on port %v\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
