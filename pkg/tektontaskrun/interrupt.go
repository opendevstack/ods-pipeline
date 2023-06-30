package tektontaskrun

import (
	"os"
	"os/signal"
)

// cleanupOnInterrupt will execute the function cleanup if an interrupt signal is caught
func cleanupOnInterrupt(cleanup func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			cleanup()
			os.Exit(1)
		}
	}()
}
