package main

import (
	"bufio"
	"net"
	"strings"
	"unicode/utf8"
)

type Client struct {
	nickName string
	conn     net.Conn
	reader   *bufio.Reader
	server   *Server
}

func nickNameInit(client *Client) (bool, error) {
	text, err := client.reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if err != nil {
		return false, err
	}

	client.nickName = strings.TrimSpace(text)
	if client.nickName == "" || utf8.RuneCountInString(client.nickName) > 9 {
		client.conn.Write([]byte("Please, insert a valid nickName beetwen 1 characher and 9\n"))
		return false, nil
	}

	reply := make(chan bool)

	client.server.nickCheck <- NickRequest{nickName: client.nickName, resultCh: reply}
	if <-reply == true {
		client.conn.Write([]byte("This name has already been taken, retry with different one.\n"))
		return false, nil
	}

	return true, nil
}

func insertPassword(client *Client) (bool, error) {

	var text, err = client.reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	text = strings.TrimSpace(text)
	if text != client.server.password {
		return false, nil
	}

	return true, nil
}

func doAuthentication(client *Client) (bool, error) {
	var attempts int8 = 0
	var passwordOk bool = false
	var err error = nil
	for attempts < 3 {
		passwordOk, err = insertPassword(client)
		if err != nil {
			return false, err
		}

		if passwordOk != true {
			client.conn.Write([]byte("Invalid password insert\n"))
		} else {
			return true, nil
		}
		attempts++
	}

	if passwordOk == false {
		client.conn.Write([]byte("Invalid limit password reached, disconnecting...\n"))
		return false, nil
	}

	return true, nil
}

func initClient(client *Client, server *Server, conn net.Conn) {
	client.server = server
	client.conn = conn
	client.reader = bufio.NewReader(client.conn)
}

func handleConnection(conn net.Conn, server *Server) {
	defer conn.Close()

	client := &Client{}
	initClient(client, server, conn)

	client.conn.Write([]byte("insert server password for access\n"))
	isAuth, err := doAuthentication(client)
	if err != nil || isAuth != true {
		return
	}

	client.conn.Write([]byte("Insert your nickName min:1 characher, max: 9 charachers:\n"))
	for {
		var ok, err = nickNameInit(client)
		if err != nil {
			return
		}
		if ok != true {
			continue
		}
		break
	}

	server.registered <- client

	for {

		text, err := client.reader.ReadString('\n')
		if err != nil {
			server.unregistered <- client
			return
		}

		text = strings.TrimSpace(text)
		if text != "" {
			server.messageCh <- &Message{
				sender:  client,
				message: text,
			}
		}
	}

}
