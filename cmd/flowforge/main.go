package main

import (
	"context"
	"os"

	"flowforge/internal/command"
)

func main() {
	if err := command.NewRootCmd().ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}
