package main

import (
	"log"
	"strings"
	"unicode/utf8"
)

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

func sendToChannel(server *Server, messageTx *Message) {

	channelName := strings.TrimPrefix(messageTx.target, "#")

	channel := server.rooms[channelName]
	if channel == nil {
		messageTx.sender.conn.Write([]byte("Channel not found, veryfi channel name\n"))
		return
	}

	if !channel.clients[messageTx.sender] {
		messageTx.sender.conn.Write([]byte("You are not channel member! Please join the channel before sending message!\n"))
		return
	}

	finalMessage := channelName + ": " + messageTx.sender.nickName + ": " + messageTx.message + "\n"
	for currentClient := range channel.clients {
		if currentClient.nickName == messageTx.sender.nickName {
			continue
		}
		currentClient.conn.Write([]byte(finalMessage))
	}
}

func addClient(client *Client, server *Server) {
	log.Println("New client arrived!")
	client.conn.Write([]byte("Welcome to the go irc server made by Giovanni Pirozzi!\n"))
	server.clients[client] = true
}

func removeClientFromChannels(client *Client, server *Server) {
	for channelName, channel := range server.rooms {
		if channel.clients[client] {
			sendToChannel(server, &Message{
				sender:  client,
				target:  channelName,
				message: "left the channel",
			})
			delete(channel.clients, client)
		}

		if len(channel.clients) == 0 {
			delete(server.rooms, channelName)
			log.Println("Channel " + channelName + " is empty, eliminated...")
		}
	}
}

func removeClient(client *Client, server *Server) {
	log.Println("Client is leaving!")
	finalMessage := client.nickName + " is leaving the server...bye\n"
	removeClientFromChannels(client, server)
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
			removeClientFromChannels(client, server)
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

func createChannel(server *Server, joinReq *JoinReq) {
	newRoom := &Room{
		name:    joinReq.roomName,
		clients: make(map[*Client]bool),
	}
	newRoom.clients[joinReq.client] = true
	server.rooms[joinReq.roomName] = newRoom
	log.Println("New room created, named: ", joinReq.roomName)
	welcomeMessage := "You are the first member of this channel named " + newRoom.name + "!\nWelcome " + joinReq.client.nickName + "\n"
	joinReq.client.conn.Write([]byte(welcomeMessage))
}

func handleJoinChannel(server *Server, joinReq *JoinReq) {

	trimmedName := strings.TrimSpace(joinReq.roomName)
	client := joinReq.client

	if trimmedName == "" || utf8.RuneCountInString(trimmedName) > 10 {
		client.conn.Write([]byte("Invalid channel name, min 1 character max 10\n"))
		return
	}

	channel := server.rooms[joinReq.roomName]

	if channel == nil {
		if strings.Contains(trimmedName, "#") {
			client.conn.Write([]byte("Invalid channel name, contains #\n"))
			return
		}
		createChannel(server, joinReq)
		return
	}

	channel.clients[client] = true
	finalMessage := "Welcome in the channel " + channel.name + "\n"
	client.conn.Write([]byte(finalMessage))
	sendToChannel(server, &Message{
		sender:  client,
		target:  channel.name,
		message: "is joined",
	})
}

func handleLeaveChannel(server *Server, leaveReq *LeaveReq) {
	trimmedName := strings.TrimPrefix(leaveReq.roomName, "#")
	trimmedName = strings.TrimSpace(trimmedName)
	client := leaveReq.client

	if trimmedName == "" || utf8.RuneCountInString(trimmedName) > 10 {
		client.conn.Write([]byte("Channel not found\n"))
		return
	}

	channel := server.rooms[trimmedName]
	if channel == nil {
		client.conn.Write([]byte("Channel not found\n"))
		return
	}

	if !channel.clients[client] {
		client.conn.Write([]byte("You are not in this channel!\n"))
		return
	}

	sendToChannel(server, &Message{
		sender:  client,
		target:  channel.name,
		message: "left the channel",
	})
	client.conn.Write([]byte("you left the channel: " + channel.name + "\n"))

	delete(channel.clients, client)

	if len(channel.clients) == 0 {
		delete(server.rooms, channel.name)
		log.Println("Channel " + channel.name + " is empty, eliminated...")
	}

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
			} else if strings.HasPrefix(messageTx.target, "#") {
				sendToChannel(server, messageTx)
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

		case joinReq := <-server.joinCh:
			handleJoinChannel(server, joinReq)

		case leaveReq := <-server.leaveCh:
			handleLeaveChannel(server, leaveReq)
		}
	}
}
