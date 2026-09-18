package main

import (
	"context"
	"log"
	"net"
	"os/signal"
	"syscall"
)

func shutdownSignal(server *Server, listener net.Listener) {
	ctx, stop := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("Shutdown signal reached, disconnecting all clients")
	server.shutDown <- true

	listener.Close()
}
