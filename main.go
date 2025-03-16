package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/zercos/oauth-tower/internal/api"
)

func main() {
	godotenv.Load()
	serv := api.CreateServer()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		serv.Logger.Info("Starting server")
		if err := serv.Start(getServerAddr()); err != nil && err != http.ErrServerClosed {
			serv.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	serv.Logger.Info("Got interrupt signal, shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := serv.Shutdown(ctx); err != nil {
		serv.Logger.Fatal(err)
	}
}

func getServerAddr() string {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000"
	}
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	return host + ":" + port
}
