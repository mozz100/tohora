package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"time"
)

var subProc *exec.Cmd
var ctx context.Context
var cancel context.CancelFunc
var cmd string

type instruction struct {
	instruction string
	parameter   string
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("Provide a command")
	}
	cmd = os.Args[1]
	log.Printf("Command is '%v'", cmd)

	var channel = make(chan instruction)
	go controller(channel)

	channel <- instruction{"start", "3600"}

	// time.AfterFunc(1*time.Second, func() {
	// 	channel <- instruction{"stop", ""}
	// })
	time.AfterFunc(2*time.Second, func() {
		channel <- instruction{"start", "3600"}
	})
	time.AfterFunc(4*time.Second, func() {
		channel <- instruction{"stop", ""}
	})

	for {
		//fmt.Print(".")
		time.Sleep(200 * time.Millisecond)
	}
}

func act(action instruction) {
	log.Printf("Action: %v", action)

	switch action.instruction {
	case "start":
		// Start the subprocess
		log.Print("Starting")
		startSubprocess(action.parameter)
		log.Print("Started")
	case "stop":
		// Stop the subprocess
		log.Print("Stopping")
		stopSubprocess()
		log.Print("Stopped")
	}
}

func startSubprocess(parameter string) {
	ctx, cancel = context.WithCancel(context.Background())
	log.Printf("Running '%v %v'", cmd, parameter)
	subProc = exec.CommandContext(ctx, cmd, parameter)
	defer cancel()

	if err := subProc.Start(); err != nil {
		log.Fatal(err)
	}
}

func stopSubprocess() {
	cancel()
	subProc.Wait() // we don't care about subprocess errors for now
}

func controller(c chan instruction) {
	for {
		action := <-c
		act(action)
	}
}
