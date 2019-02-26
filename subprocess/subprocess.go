package subprocess

import (
	"context"
	"log"
	"os/exec"
)

var subProc *exec.Cmd
var ctx context.Context
var cancel context.CancelFunc

// Instruction - what to do, and a Parameter
type Instruction struct {
	Instruction string
	Parameter   string
}

// Subprocess returns the channel for sending Instructions on
func Subprocess(cmd string) *chan Instruction {
	var channel = make(chan Instruction)
	go controller(cmd, channel)

	return &channel
}

func act(cmd string, action Instruction) {
	log.Printf("Action: %v", action)

	switch action.Instruction {
	case "start":
		// Start the subprocess
		log.Print("Starting")
		startSubprocess(cmd, action.Parameter)
		log.Print("Started")
	case "stop":
		// Stop the subprocess
		log.Print("Stopping")
		stopSubprocess()
		log.Print("Stopped")
	}
}

func startSubprocess(cmd string, Parameter string) {
	ctx, cancel = context.WithCancel(context.Background())
	log.Printf("Running '%v %v'", cmd, Parameter)
	subProc = exec.CommandContext(ctx, cmd, Parameter)
	defer cancel()

	if err := subProc.Start(); err != nil {
		log.Fatal(err)
	}
}

func stopSubprocess() {
	cancel()
	subProc.Wait() // we don't care about subprocess errors for now
}

func controller(cmd string, c chan Instruction) {
	for {
		action := <-c
		act(cmd, action)
	}
}
