package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"time"
)

var subProc *exec.Cmd

func main() {
	var channel = make(chan string)

	go controller(channel)

	channel <- "start"

	time.AfterFunc(1*time.Second, func() {
		channel <- "stop"
	})
	time.AfterFunc(2*time.Second, func() {
		channel <- "start"
	})
	time.AfterFunc(4*time.Second, func() {
		channel <- "stop"
	})

	for {
		fmt.Print(".")
		time.Sleep(200 * time.Millisecond)
	}
}

var stdout io.ReadCloser
var ctx context.Context
var cancel context.CancelFunc

func act(action string) {
	fmt.Print("action: ", action)
	var err1 error

	switch action {
	case "start":
		ctx, cancel = context.WithCancel(context.Background())
		subProc = exec.CommandContext(ctx, "bash", "-c", "while true; do echo -n 'parp'; sleep 0.2; done")
		defer cancel()
		if stdout, err1 = subProc.StdoutPipe(); err1 != nil {
			log.Fatal(err1)
		}

		if err := subProc.Start(); err != nil {
			log.Fatal(err)
		}
		fmt.Print("started")
	case "stop":
		fmt.Print("stopping")

		buffer := make([]byte, 100)
		fmt.Print(string(buffer))
		stdout.Read(buffer)
		stdout.Close()
		cancel()
		err := subProc.Wait()
		log.Printf("Command finished with error: %v", err)
		fmt.Print("stopped")
	}
}

func controller(c chan string) {
	for {
		action := <-c
		act(action)
	}
}
