package main

import (
	"fmt"
	"log"

	"github.com/mellonnen/chronograph/ui"
)

func main() {
	program := ui.New("chronograph.db")
	if err := program.Start(); err != nil {
		log.Fatal(fmt.Errorf("initializing UI: %w", err))
	}
}
