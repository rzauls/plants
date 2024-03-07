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
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		cancel()
		os.Exit(0)
	}()

	if err := httpd.Run(ctx, args, getenv, stdin, stdout, stderr); err != nil {
		return err
	}

	return nil
}

func main() {
	ctx := context.Background()
	// NOTE: in main() we use all the os defaults
	if err := run(ctx, os.Args, os.Getenv, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
