package main

import (
	"fmt"
	"bufio"
	"net"
)

type ClientMap map[*Client]bool

type Client struct {
	username	string
	conn		*net.Conn
	room		*Room
	// Receive/Send channel/s
	Rcv		chan string
}

func CreateClient(username string, room *Room, conn *net.Conn) *Client {
	return &Client{username, conn, room, make(chan string)}
}

func (c *Client) WriteMessage() {
	// Infinite for loop that's always listening for input
	// That input is sent to the broadcaster of the client's room
	for {
		for _, msg := range <-c.Rcv {
			fmt.Fprint(*c.conn, string(msg))
		}
		fmt.Println()
	}
}

func (c *Client) ReadMessage() {
	for {
		input := bufio.NewScanner(*c.conn)
		for input.Scan() {
			fmt.Fprint(*c.conn, "> ")
			// Send input.Text() to the room's broadcaster
			c.room.Broadcaster.messages <- c.username + ": " + input.Text() + "\n"
		}
	}
}
