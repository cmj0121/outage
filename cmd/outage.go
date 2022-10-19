package main

import (
	"os"

	"github.com/cmj0121/outage"
)

func main() {
	agent := outage.New()

	if err := agent.Run(); err != nil {
		os.Exit(1)
		return
	}
}
