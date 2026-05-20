package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"ticketing-system/internal/server"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func gracefulShutdown(fiberServer *server.FiberServer, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal from hardware system triggers.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force closure...")
	stop() // Allow Ctrl+C to force absolute execution cutoff

	// The context is used to inform the server it has 5 seconds to finish its remaining database transactions.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := fiberServer.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server processing loops exited.")

	// Notify the main goroutine thread execution that the shutdown task is finished
	done <- true
}

func main() {
	// Spin up configuration initialization, DB layers wiring, and routing setups
	srv := server.New()

	// Create a done communication channel to track termination completeness status
	done := make(chan bool, 1)

	// Execute listening background channel processing threads runtime loop
	go func() {
		portStr := os.Getenv("PORT")
		if portStr == "" {
			portStr = "8080" // Default port if variable missing
		}
		
		port, _ := strconv.Atoi(portStr)
		log.Printf("HTTP Web Server running locally on port :%d", port)
		
		err := srv.Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			log.Fatalf("http server crash execution fault: %s", err)
		}
	}()

	// Activate operational operating-system signals surveillance monitoring thread loop
	go gracefulShutdown(srv, done)

	// Block processing until completion notification is signaled from channel
	<-done
	log.Println("Graceful shutdown sequence completed. Application safe to turn off.")
}
