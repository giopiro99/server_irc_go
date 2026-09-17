package main

func runBroadCaster(server *Server) {

	for {
		select {
		case client := <-server.registered:
			println("New client arrived!")
			client.conn.Write([]byte("Welcome to the go irc server made by Giovanni Pirozzi!\n"))
			server.clients[client] = true

		case client := <-server.unregistered:
			println("Client disconnected from the server")
			if server.clients[client] {
				delete(server.clients, client)
			}

		case messageTx := <-server.messageCh:
			for client := range server.clients {
				if client == messageTx.sender {
					continue
				}
				client.conn.Write([]byte(messageTx.message))
			}
		}
	}
}
