package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
	"strings"
)

type Client struct {
	conn     net.Conn
	nickname string
}

var clients = make(map[net.Conn]*Client)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error creating listener:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Chat server started. Listening on port 8080.")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	client := &Client{
		conn: conn,
	}

	clients[conn] = client

	fmt.Println("New client connected:", conn.RemoteAddr().String())

	sendMessageToClient(client, "Welcome to the chat! Please enter your nickname:")

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Text()

		if client.nickname == "" {
			client.nickname = message
			sendMessageToAllClients(client, fmt.Sprintf("%s joined the chat!", client.nickname))
		} else {
			if len(message) > 200 {
				sendErrorMessageToClient(client, "Message is too long. Please shorten it.")
				continue
			}
			sendMessageToAllClients(client, fmt.Sprintf("%s: %s", client.nickname, message))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from client:", err)
	}

	conn.Close()

	delete(clients, conn)

	sendMessageToAllClients(client, fmt.Sprintf("%s left the chat.", client.nickname))

	fmt.Println("Client disconnected:", conn.RemoteAddr().String())

}

func sendMessageToClient(client *Client, message string) {
	client.conn.Write([]byte(message + "\n"))
}

func userList() string {
	var userList strings.Builder
	for _, client := range clients {
		if client.nickname != "" {
			userList.WriteString(client.nickname)
			userList.WriteString(", ")
		}
	}
	if userList.Len() > 2 {
		userListString := userList.String()
		userListString = userListString[:userList.Len()-2] 
		return userListString
	}
	return ""
}

func sendMessageToAllClients(sender *Client, message string) {
	timestamp := time.Now().Format("15:04:05")
	messageWithTimestamp := fmt.Sprintf("%s [%s]", message, timestamp)

	userListMsg := fmt.Sprintf("Users online: %s", userList())
	if userListMsg != "Users online: " {
		messageWithTimestamp += " | " + userListMsg
	}

	for conn, client := range clients {
		if conn != nil && client != sender {
			sendMessageToClient(client, messageWithTimestamp)
		}
	}
}

func sendErrorMessageToClient(client *Client, message string) {
	errorMsg := "Error: " + message

	client.conn.Write([]byte(errorMsg + "\n"))
}


