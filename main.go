package main

import (
	"log"
	"os"
	"os/signal"
	"plants/httpd"
)

// basic http api to test swagger annotations
func main() {
	// handle exit signals
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		log.Println("Program terminated by OS signal")
		// TODO: any additional things to do before exiting
		os.Exit(0)
	}()

	httpd.Run()
}
