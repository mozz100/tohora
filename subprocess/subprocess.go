package subprocess

import (
	"context"
	"log"
	"os/exec"
)

// Context holds together the bits
type Context struct {
	SubProc *exec.Cmd
	Ctx     context.Context
	Cancel  context.CancelFunc
}

// Instruction - what to do, and a Parameter
type Instruction struct {
	Instruction string
	Parameter   string
}

// Subprocess returns the channel for sending Instructions on
func Subprocess(cmd string, sbpctx *Context) *chan Instruction {
	var channel = make(chan Instruction)
	go controller(cmd, sbpctx, channel)

	return &channel
}

func act(cmd string, action Instruction, sbpctx *Context) {
	log.Printf("Action: %v", action)

	switch action.Instruction {
	case "start":
		// Start the subprocess
		log.Print("Starting")
		startSubprocess(cmd, action.Parameter, sbpctx)
		log.Print("Started")
	case "stop":
		// Stop the subprocess
		log.Print("Stopping")
		stopSubprocess(sbpctx)
		log.Print("Stopped")
	}
}

func startSubprocess(cmd string, Parameter string, sbpctx *Context) {
	sbpctx.Ctx, sbpctx.Cancel = context.WithCancel(context.Background())
	log.Printf("Running '%v %v'", cmd, Parameter)
	sbpctx.SubProc = exec.CommandContext(sbpctx.Ctx, cmd, Parameter)
	defer sbpctx.Cancel()

	if err := sbpctx.SubProc.Start(); err != nil {
		log.Fatal(err)
	}
}

func stopSubprocess(sbpctx *Context) {
	if sbpctx.SubProc == nil {
		return
	}
	sbpctx.Cancel()
	sbpctx.SubProc.Wait() // we don't care about subprocess errors for now
}

func controller(cmd string, sbpctx *Context, c chan Instruction) {
	for {
		action := <-c
		act(cmd, action, sbpctx)
	}
}
