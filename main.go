package main

import (
	"errors"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
)

type Message struct {
	sender  *Client
	message string
}

type NickRequest struct {
	nickName string
	resultCh chan bool
}

type Server struct {
	ipAddress string
	port      string
	password  string

	messageCh    chan *Message
	clients      map[*Client]bool
	registered   chan *Client
	unregistered chan *Client

	shutDown  chan bool
	nickCheck chan NickRequest
}

func newServer() (*Server, error) {

	server := &Server{}

	server.ipAddress = os.Getenv("IP_ADDRESS")
	if server.ipAddress == "" {
		server.ipAddress = "127.0.0.1"
	}

	server.port = os.Getenv("PORT")
	if server.port == "" {
		server.port = "8190"
	}

	server.password = os.Getenv("PASSWORD")
	if server.password == "" {
		return nil, errors.New("PASSWORD is empty")
	}

	server.messageCh = make(chan *Message)
	server.clients = make(map[*Client]bool)
	server.registered = make(chan *Client)
	server.unregistered = make(chan *Client)
	server.shutDown = make(chan bool)
	server.nickCheck = make(chan NickRequest)
	return server, nil
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No file .env, exit")
		return
	}

	server, err := newServer()
	if err != nil {
		log.Fatal("Impossible to start server: ", err)
	}

	log.Println("Configuration applyed, server starting on port: ", server.port)

	connection := net.JoinHostPort(server.ipAddress, server.port)
	listener, err := net.Listen("tcp", connection)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	defer listener.Close()

	go shutdownSignal(server, listener)
	go runBroadCaster(server)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error: ", err)
			break
		}
		go handleConnection(conn, server)
	}

}
