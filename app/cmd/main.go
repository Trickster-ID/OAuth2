package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"log"
	"oauth2/app/helper"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	helper.InitialConfig()
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "main")
	}

	if os.Args[1] == "main" {
		fiberServer := InitializeFiberServer()
		startHttpServer(fiberServer)
	}

}

func startHttpServer(f *fiber.App) {

	// Channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	mainPort, err := strconv.Atoi(os.Getenv("MAIN_PORT"))
	if err != nil {
		logrus.Fatal(err)
	}
	go func() {
		logrus.Info("Running server on port ", mainPort)
		if err := f.Listen(fmt.Sprintf(":%d", mainPort)); err != nil {
			logrus.Error(err)
			logrus.Fatal("shutting down http server")
		}
	}()

	// Block until a signal is received
	<-quit
	log.Print("Server is shutting down...")

	// Create a deadline for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := f.ShutdownWithContext(ctx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Print("Server stopped")
}
