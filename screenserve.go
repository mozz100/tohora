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
		<html>
			<head></head>
			<body>
				<form method='POST' action='/launch/'>
					<input name='url' type='text' />
					<input type='submit'/>
				</form>
			</body>
		</html>`)
		log.Println("Responded with home page")
	}

	http.HandleFunc("/launch/", launchHandler)
	http.HandleFunc("/", homeHandler)
	log.Printf("Starting listening on port %v\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
