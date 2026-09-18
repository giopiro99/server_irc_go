package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/joho/godotenv"
)

type Message struct {
	sender  *Client
	target  string
	message string
}

type NickRequest struct {
	nickName string
	resultCh chan bool
}

type JoinReq struct {
	client   *Client
	roomName string
}

type LeaveReq struct {
	client   *Client
	roomName string
}

type Room struct {
	name    string
	clients map[*Client]bool
}

type Server struct {
	ipAddress string
	port      string
	password  string

	clients map[*Client]bool
	rooms   map[string]*Room

	joinCh       chan *JoinReq
	leaveCh      chan *LeaveReq
	messageCh    chan *Message
	registered   chan *Client
	unregistered chan *Client
	shutDown     chan bool
	nickCheck    chan NickRequest
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
	server.rooms = make(map[string]*Room)
	server.joinCh = make(chan *JoinReq)
	server.leaveCh = make(chan *LeaveReq)

	return server, nil
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("No file .env, exit")
		return
	}

	server, err := newServer()
	if err != nil {
		log.Fatalf("Impossible to start server: %v", err)
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

	go func() {
		log.Println("Server pprof in ascolto su http://localhost:6060/debug/pprof/")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error: ", err)
			break
		}
		go handleConnection(conn, server)
	}

}
