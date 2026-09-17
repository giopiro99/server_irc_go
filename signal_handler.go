package main

import (
	"context"
	"net"
	"os/signal"
	"syscall"
)

func shutdownSignal(server *Server, listener net.Listener) {
	ctx, stop := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	println("Shutdown signal reached, disconnecting all clients\n")
	server.shutDown <- true

	listener.Close()
}
