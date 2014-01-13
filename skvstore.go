package skvstore

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

var (
	listener net.Listener // Global copy of listener

	dict           = map[string]string{} // The key-value dictionary
	port           = "6428"              // The server port
	server_running = false               // State of the server
)

// Stops the server if it was running
func stopServer() {
	if server_running {
		server_running = false
		listener.Close()
	}
}

// Stats the SKVStore server
func startServer() {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Panicln("Could not listen on port " + port)
	}
	log.Println("[+] SKVStore Server Listerning on port " + port)

	server_running = true
	listener = lis
	for {
		conn, err := listener.Accept()
		if err != nil && server_running {
			panic("Could not accept connection. Error: " + err.Error())
		} else if !server_running {
			log.Println("[-] Stopped SKVStore Server.")
			return
		}
		log.Println("[@] SKVStore Server accepted a connection from " + conn.RemoteAddr().String())

		go handle(conn)
	}
}

// Handles an incoming connection from a client
func handle(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		data, err := reader.ReadString('\n')
		if err == io.EOF {
			log.Println("[X] Client (" + conn.RemoteAddr().String() + ") has been disconnected.")
			break
		} else if err != nil {
			log.Println("[!] Error reading data from client (" + conn.RemoteAddr().String() + ") :: " + err.Error())
			return
		}
		data = data[:len(data)-2]
		log.Println("[.] Received data from client (" + conn.RemoteAddr().String() + ") :: " + data)

		parts := strings.Split(data, " ")
		command := parts[0]
		message := ""

		switch command {
		case "GET":
			if len(parts) <= 1 {
				message = "NO KEY SPECIFIED!"
			} else {
				key := parts[1]
				val, found := dict[key]
				if found {
					message = "FOUND " + key + " => " + val
				} else {
					message = "KEY " + key + " NOT FOUND!"
				}
			}

		case "SET":
			if len(parts) <= 1 {
				message = "NO KEY SPECIFIED!"
			} else {
				key := parts[1]
				if len(parts) <= 2 {
					dict[key] = parts[2]
					message = "STORED " + key + " => " + dict[key]
				} else {
					message = "NO VALUE PROVIDED FOR " + key + "!"
				}
			}

		case "CLEAR":
			dict = make(map[string]string)
			message = "ALL ENTRIES ERASED"

		case "PRINT":
			message = string(len(dict))
			for key, value := range dict {
				message += "\r\n" + key + " => " + value
			}

		default:
			message = "UNKNOWN COMMAND -- " + command
		}

		log.Println("[.] Sending message to client (" + conn.RemoteAddr().String() + ") :: " + message)
		conn.Write([]byte(message + "\r\n"))
	}
}
