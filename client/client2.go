package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type client struct {
}

func NewClient() *client {
	return &client{}
}

func (c *client) Start(address string) {
	connection,err := net.Dial("tcp", address)

	if err!=nil {
		fmt.Printf("Error while dialling the server %s: ", connection.RemoteAddr())
		fmt.Println("")
	}

	go listenServerLoop(connection)
	sendMessageLoop(connection)
}

func sendMessageLoop(connection net.Conn) {
	fmt.Println("Send Message Loop has started")

	for{

		reader := bufio.NewReader(os.Stdin)
		readLine,err := reader.ReadString('\n')

		if err!=nil {
			fmt.Printf("Error while reading from console: %s", err.Error())
		}

		command := strings.TrimSpace(readLine)
		// fmt.Printf("The line read is: %s", command)

		if strings.ToLower(command) == "reconnect" {
			newConnection,_ := net.Dial("tcp", ":3000")
			connection = newConnection
			fmt.Println("Reconnection is done")
		}

		connection.Write([]byte(readLine))
	}
}

func listenServerLoop(connection net.Conn) {
	// listener,err1 := net.Listen("tcp", connection.LocalAddr().String())
	var buf []byte = make([]byte, 2048)
	fmt.Println("Hello")
	for {
		// if err1 != nil {
		// 	fmt.Printf("Error while creating listener to server: %s", err1.Error())
		// 	fmt.Println(" ")
		// }
	
		// conn,err2 := listener.Accept();
	
		// if err2!=nil {
		// 	fmt.Printf("Error while accepting connection to server: %s", err2.Error())
		// 	fmt.Println(" ")
		// }

		fmt.Println("Hello")
		n,_ := connection.Read(buf)

		fmt.Printf("What was read from the server was: %s", string(buf[:n]))
		fmt.Println(" ")
	}
}


func main() {
	client := NewClient()
	client.Start("127.0.0.1:3000")
}
