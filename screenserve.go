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
		param := r.FormValue("param")
		log.Printf("Launch page with param = '%v'\n", param)
		if sbpctx.SubProc != nil || param == "" {
			ch <- subprocess.Instruction{
				Instruction: "stop",
				Parameter:   "",
			}
		}
		if param != "" {
			ch <- subprocess.Instruction{
				Instruction: "start",
				Parameter:   param,
			}
		}

		log.Println("Redirecting to home page")
		http.Redirect(w, r, "/", http.StatusFound)
	}

	homeHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
		<html>
			<head></head>
			<body>
				<form method='POST' action='/launch/'>
					<input name='param' type='text' />
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
