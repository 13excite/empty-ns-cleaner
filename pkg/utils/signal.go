package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
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
			zap.S().Info("received interrupt...")
			return
		case <-ctx.Done():
			zap.S().Info("shutting down by context...")
			return
		}
	}
}
