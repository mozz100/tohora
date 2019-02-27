package subprocess

import (
	"context"
	"log"
	"os/exec"
)

// Context holds together the bits
type Context struct {
	SubProc   *exec.Cmd
	Ctx       context.Context
	Cancel    context.CancelFunc
	Channel   chan string
	Parameter string
}

// Subprocess returns the channel for sending instructions on
func Subprocess(cmd string) *Context {
	sbpctx := Context{}
	sbpctx.Channel = make(chan string)
	go controller(cmd, &sbpctx)

	return &sbpctx
}

func startSubprocess(cmd string, Parameter string, sbpctx *Context) {
	sbpctx.Ctx, sbpctx.Cancel = context.WithCancel(context.Background())
	log.Printf("Running '%v %v'", cmd, Parameter)
	sbpctx.SubProc = exec.CommandContext(sbpctx.Ctx, cmd, Parameter)

	if err := sbpctx.SubProc.Start(); err != nil {
		log.Fatal(err)
	}
}

func stopSubprocess(sbpctx *Context) {
	sbpctx.Cancel()
	sbpctx.SubProc.Wait() // we don't care about subprocess errors for now
}

func controller(cmd string, sbpctx *Context) {
	for {
		action := <-sbpctx.Channel
		log.Printf("Action: %v", action)
		sbpctx.Parameter = action
		switch action {
		case "":
			// Stop the subprocess
			log.Print("Stopping")
			stopSubprocess(sbpctx)
			log.Print("Stopped")
		default:
			// Start the subprocess
			log.Print("Starting")
			startSubprocess(cmd, action, sbpctx)
			log.Print("Started")
		}
	}
}
