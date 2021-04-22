package main

import (
	"os"
	"os/signal"

	"github.com/mwennrich/partition-watchdog/cmd"
)

func main() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc)
	go func() {
		<-sigc
	}()
	cmd.Execute()
}
