package subprocess

import (
	"context"
	"log"
	"os/exec"
)

// Context contains everything needed - the underlying Cmd and the
// mechanism to terminate it.
type Context struct {
	Command   string
	Parameter string
	History   []string

	subProc *exec.Cmd
	ctx     context.Context
	cancel  context.CancelFunc
}

// GetSubprocess sets up a context and returns it
func GetSubprocess(cmd string) *Context {
	sbpctx := Context{}
	sbpctx.Command = cmd
	sbpctx.History = make([]string, 10)
	return &sbpctx
}

// StartWith starts a subprocess using the specified Command and a Parameter
func (sbpctx *Context) StartWith(Parameter string) {
	sbpctx.ctx, sbpctx.cancel = context.WithCancel(context.Background())
	log.Printf("Running '%v %v'", sbpctx.Command, Parameter)
	sbpctx.subProc = exec.CommandContext(sbpctx.ctx, sbpctx.Command, Parameter)

	if err := sbpctx.subProc.Start(); err != nil {
		log.Fatal(err)
	}
	sbpctx.Parameter = Parameter
	if Parameter != sbpctx.History[len(sbpctx.History)-1] {
		sbpctx.History = sbpctx.History[1:]
		sbpctx.History = append(sbpctx.History, Parameter)
	}
}

// Stop is used to stop the subprocess and set Parameter to ""
func (sbpctx *Context) Stop() {
	if sbpctx.subProc == nil {
		return
	}
	log.Println("Stopping")
	sbpctx.cancel()
	sbpctx.subProc.Wait() // we don't care about subprocess errors for now
	sbpctx.Parameter = ""
}
