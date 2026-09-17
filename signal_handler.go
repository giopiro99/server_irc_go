package main

import (
	"context"
	"os/signal"
	"syscall"
)

func shutdownSignal(server *Server) {
	ctx, stop := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	println("Shutdown signal reached, disconnecting all clients\n")
	server.shutDown <- true
}
