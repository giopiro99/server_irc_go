package main

import "log"

func sendBroadCast(server *Server, messageTx *Message) {
	for client := range server.clients {
		if client == messageTx.sender {
			continue
		}
		var finalMessage string = messageTx.sender.nickName + ": " + messageTx.message + "\n"
		client.conn.Write([]byte(finalMessage))
	}
}

func sendPrivMsg(server *Server, messageTx *Message) {
	for client := range server.clients {
		if client == messageTx.sender {
			continue
		}
		if client.nickName == messageTx.target {
			var finalMessage string = messageTx.sender.nickName + ": " + messageTx.message + "\n"
			client.conn.Write([]byte(finalMessage))
			break
		}
	}
}

func addClient(client *Client, server *Server) {
	log.Println("New client arrived!")
	client.conn.Write([]byte("Welcome to the go irc server made by Giovanni Pirozzi!\n"))
	server.clients[client] = true
}

func removeClient(client *Client, server *Server) {
	log.Println("Client is leaving!")
	finalMessage := client.nickName + " is leaving the server...bye\n"
	for currentClient := range server.clients {
		if client == currentClient {
			continue
		}
		currentClient.conn.Write([]byte(finalMessage))
	}
	delete(server.clients, client)
}

func removeAllClients(server *Server) {
	for client := range server.clients {
		if server.clients[client] {
			client.conn.Write([]byte("Error: Server is in shutdown, disconnecting..\n"))
			delete(server.clients, client)
			client.conn.Close()
		}
	}
}

func checkNickName(server *Server, req NickRequest) bool {
	for client := range server.clients {
		if client.nickName == req.nickName {
			return true
		}
	}
	return false
}

func runBroadCaster(server *Server) {

	for {
		select {
		case client := <-server.registered:
			addClient(client, server)
		case client := <-server.unregistered:
			if server.clients[client] {
				removeClient(client, server)
			}
		case messageTx := <-server.messageCh:
			if messageTx.target == "" {
				sendBroadCast(server, messageTx)
			} else {
				sendPrivMsg(server, messageTx)
			}
		case shutDown := <-server.shutDown:
			if shutDown == true {
				removeAllClients(server)
				return
			}
		case req := <-server.nickCheck:
			req.resultCh <- checkNickName(server, req)
		}
	}
}
