package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mozz100/screenserve/handlers"
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

	http.HandleFunc("/", handlers.HomeHandler(sbpctx))
	http.HandleFunc("/launch/", handlers.LaunchHandler(sbpctx))

	log.Printf("Starting listening on port %v\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
