package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Server struct {
	connectionList map[net.Addr]net.Conn
	database map[string][]byte
	quitch chan struct{}
	reqChannel chan string
	getResultChannel chan []byte
	listener net.Listener
}

func createServer() *Server {
	return &Server{
		connectionList: make(map[net.Addr]net.Conn),
		database: make(map[string][]byte),
		quitch: make(chan struct{}),
		reqChannel: make(chan string),
		getResultChannel : make(chan []byte),
		listener: nil,
	}
}

func (s *Server) Start() {
	go s.dbOpsListener()
	go s.listenForInterrupt()

	listener, err := net.Listen("tcp", "127.0.0.1:3000");
	s.listener = listener
	fmt.Println("The server has started")

	if err!=nil {
		fmt.Printf("Error occured while receiving connection from: %s. The error is: %s\n", listener.Addr(), err)
	}

	for {
		s.connectionAcceptLoop(listener)
	}
}

func (s *Server) connectionAcceptLoop(listener net.Listener) {
	connection, err := listener.Accept();
	s.connectionList[connection.RemoteAddr()] = connection
	fmt.Printf("The connection has been successfully made with port: %s\n", connection.RemoteAddr())

	if nil!=err {
		fmt.Printf("Error occured while accepting connection from: %s. The error is: %s\n", listener.Addr(), err)
	}

	go s.dataReadLoop(connection)
}

func (s *Server) dataReadLoop(connection net.Conn) {
	var byteStream []byte = make([]byte, 2048);

	for {
		n, err := connection.Read(byteStream)
		// We can use connection.ReadString function as well
		// connection.ReadString and other methods may send back delimeters and spaces.

		if err!=nil {
			fmt.Println("Error occured while reading data from the stream")
			break;
		}

		fmt.Printf("The data received from %s is: %s\n", connection.RemoteAddr(), byteStream[:n])

		commandString := strings.TrimSpace(string(byteStream[:n]))
		command := strings.Split(commandString, ",")[0]

		if command == "put" {
			s.handlePut(commandString)

		} else if command == "get" {
			s.handleGet(commandString) 
		}
	}
}

func (s *Server) Count() (int) {
	return len(s.connectionList)
}

func (s *Server) handlePut(command string) {
	s.reqChannel <- command
}

func (s *Server) handleGet(command string) {
	s.reqChannel <- command
	result := <-s.getResultChannel
	for _,clientConnection := range s.connectionList {
		clientConnection.Write(result)
		fmt.Printf("Writing to server: %s\n", clientConnection.RemoteAddr().String())
	}
}

func (s *Server) dbOpsListener() {
	for {
		commandString := <-s.reqChannel
		command := strings.Split(commandString, ",")

		if command[0] == "put" {
			s.database[command[1]] = []byte(command[2])
		} else if command[0] == "get" {
			result := s.database[command[1]]
			s.getResultChannel <- result
		}
	}
}

func (s *Server) listenForInterrupt() {
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)

	<-interruptChannel
	s.listener.Close()
}

func main() {
	server := createServer()
	fmt.Println("The server has been created")
	server.Start()
}