package main

import (
	"log"
	"os"
	"time"

	subprocess "github.com/mozz100/screenserve/subprocess"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("Provide a command")
	}
	cmd := os.Args[1]
	log.Printf("Command is '%v'", cmd)

	ch := *subprocess.Subprocess(cmd)
	log.Print(ch)

	ch <- subprocess.Instruction{
		Instruction: "start",
		Parameter:   "3600",
	}

	time.Sleep(time.Second * 10)
	ch <- subprocess.Instruction{
		Instruction: "stop",
		Parameter:   "",
	}
	time.Sleep(time.Second * 10)
}
