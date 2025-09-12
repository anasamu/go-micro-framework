package main

import (
	"fmt"
	"os"

	"github.com/anasamu/go-micro-framework/cmd/microframework/commands"
)

func main() {
	if err := commands.GetRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
