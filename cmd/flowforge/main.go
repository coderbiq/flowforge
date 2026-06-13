package main

import (
	"context"
	"fmt"
	"os"

	"flowforge/internal/command"
)

func main() {
	if err := command.NewRootCmd().ExecuteContext(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
