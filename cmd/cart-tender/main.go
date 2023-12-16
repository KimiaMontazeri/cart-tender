package main

import (
	"os"

	"github.com/KimiaMontazeri/cart-tender/internal/cart-tender/cmd"
)

const (
	exitFailure = 1
)

func main() {
	root := cmd.NewRootCommand()

	if root != nil {
		if err := root.Execute(); err != nil {
			os.Exit(exitFailure)
		}
	}
}
