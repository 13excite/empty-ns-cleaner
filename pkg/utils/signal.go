package utils

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func WaitForShutdown(ctx context.Context) {
	// Handle signals
	signalChannel := make(chan os.Signal, 1)
	// Stop flag will indicate if Ctrl-C/Interrupt has been sent to the process
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-signalChannel:
			// wait for system signal && initialize termination in that case
			log.Printf("received interrupt...")
			return
		case <-ctx.Done():
			log.Printf("shutting down by context...")
			return
		}
	}
}
