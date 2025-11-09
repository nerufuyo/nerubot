package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	version = "3.0.0"
	name    = "NeruBot"
)

func main() {
	fmt.Printf("=== %s v%s (Golang Edition) ===\n", name, version)
	fmt.Println("Starting Discord bot...")

	// TODO: Initialize configuration
	// TODO: Setup logger
	// TODO: Initialize Discord session
	// TODO: Start bot

	log.Println("Bot is running. Press CTRL+C to exit.")

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("\nShutting down gracefully...")
	// TODO: Cleanup resources
	fmt.Println("Goodbye!")
}
