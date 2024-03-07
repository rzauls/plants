package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"plants/httpd"
)

// run - starts the http daemon with a cancellable context
// run is just a wrapper so we can always return errors
// its also useful for testing the service
func run(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	return httpd.Run(ctx)
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args, os.Getenv, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
