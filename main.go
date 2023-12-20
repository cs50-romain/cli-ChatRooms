package main

import (
	"fmt"
	"log"
	"net"
)

// Mostly used to close client connection if even used
type Manager struct {
	clients	ClientMap
}

func handleConn(conn *net.Conn) {
	var username string
	var roomname string

	var room *Room
	fmt.Fprint(*conn, "Enter your username: ")
	_, err := fmt.Fscan(*conn, &username)
	fmt.Fprint(*conn, "Enter a room name: ")
	_, err = fmt.Fscan(*conn, &roomname)

	if err != nil {
		log.Print("[ERROR] Error creating user or room ->", err)
	}

	// Create a new client
	client := CreateClient(username, nil, conn)

	if _, ok := RoomList[roomname]; !ok {
		room = NewRoom(roomname)
		room.Broadcaster.messages = make(chan string)
		go room.Broadcast()
		RoomList[roomname] = room
		room.Clients = append(room.Clients, client)
	} else {
		room = RoomList[roomname]
		room.Clients = append(room.Clients, client)
	}

	client.room = room

	// Goroutine to start WriteMessage and ReadMessage
	go client.ReadMessage()
	go client.WriteMessage()

	select{}
	// Might have to add a select statement so that the handleConn doesn't return and so the connection stays alive
}

func main() {
	RoomList = make(map[string]*Room)
	listener, err := net.Listen("tcp",":8080")
	if err != nil {
		log.Print("[ERROR] error with listener -> ", err)
	}

	fmt.Println("Listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print("[ERROR] Error connecting ->", err)
		}
		log.Println("[INFO] Accepted connection from:", conn.RemoteAddr())

		go handleConn(&conn)
	}
}
