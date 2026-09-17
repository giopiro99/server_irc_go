package main

import "log"

func runBroadCaster(server *Server) {

	for {
		select {
		case client := <-server.registered:
			log.Println("New client arrived!")
			client.conn.Write([]byte("Welcome to the go irc server made by Giovanni Pirozzi!\n"))
			server.clients[client] = true

		case client := <-server.unregistered:
			if server.clients[client] {
				finalMessage := client.nickName + " is leaving the server...bye\n"
				for currentClient := range server.clients {
					if client == currentClient {
						continue
					}
					currentClient.conn.Write([]byte(finalMessage))
				}
				delete(server.clients, client)
			}

		case messageTx := <-server.messageCh:
			for client := range server.clients {
				if client == messageTx.sender {
					continue
				}
				var finalMessage string = messageTx.sender.nickName + ": " + messageTx.message + "\n"
				client.conn.Write([]byte(finalMessage))
			}
		case shutDown := <-server.shutDown:
			if shutDown == true {
				for client := range server.clients {
					if server.clients[client] {
						client.conn.Write([]byte("Error: Server is in shutdown, disconnecting..\n"))
						delete(server.clients, client)
						client.conn.Close()
					}
				}
				return
			}
		case req := <-server.nickCheck:
			var nameExists bool = false
			for client := range server.clients {
				if client.nickName == req.nickName {
					nameExists = true
					break
				}
			}
			req.resultCh <- nameExists
		}
	}
}
