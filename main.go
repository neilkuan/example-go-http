package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	r := SetupRouter()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	idleConnsClosed := make(chan struct{})

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 5 seconds.
		quit := make(chan os.Signal, 1)
		// kill (no param) default send syscall.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-quit
		log.Println("Shutting down server...")

		// The context is used to inform the server it has 5 seconds to finish
		// the request it is currently handling
		time.Sleep(65 * time.Second)
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatal("Server forced to shutdown: ", err)
		}
		close(idleConnsClosed)
		log.Println("Server exiting")

	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

	<-idleConnsClosed

}
