package skvstore

import (
	"bufio"
	"net"
	"testing"
	"time"
)

func WriteConnection(conn net.Conn, message string) bool {
	_, err := conn.Write([]byte(message + "\r\n"))
	if err != nil {
		conn.Close()
		return false
	}
	return true
}

func ReadConnection(conn net.Conn) (bool, string) {
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		conn.Close()
		return false, "ERROR Reading from socket!"
	}
	return true, message[:len(message)-2]
}

func LogConnectionEvent(conn net.Conn, t *testing.T, msg string) {
	t.Log(conn.LocalAddr().String() + " -- " + msg)
}

func Test_SKVStore_1(t *testing.T) {
	go startServer()
	defer stopServer()

	for !server_running {
		time.Sleep(100 * time.Millisecond)
	}

	conn, err := net.Dial("tcp", ":"+port)
	if err != nil {
		t.Error("[!] " + err.Error())
	}
	defer conn.Close()
	LogConnectionEvent(conn, t, "[@] Connected to SKVStore Server @ port "+port)

	message := "GET"
	expected := "NO KEY SPECIFIED!"

	LogConnectionEvent(conn, t, "[>] Sending message to server :: "+message)
	WriteConnection(conn, message)
	ok, reply := ReadConnection(conn)
	if ok {
		LogConnectionEvent(conn, t, "[<] Received message from server :: "+reply)
		if reply != expected {
			t.Error("[/] Was expecting :: " + expected + "  |  received :: " + reply)
		}
	} else {
		t.Error("[!] " + reply)
		return
	}
}
