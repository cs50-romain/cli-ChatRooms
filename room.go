package main

var RoomList map[string]*Room

type Broadcaster struct {
	messages chan string
}

type Room struct {
	roomname	string
	Id		int
	Clients		[]*Client
	Broadcaster	*Broadcaster
}

func (r *Room) Broadcast() {
	// Receive and send messages from and to clients of the room
	for  {
		select {
		case msg := <-r.Broadcaster.messages:
			// Broadcast message to all client in the room
			for _, client := range r.Clients {
				client.Rcv <- msg
			}
		}
	}
}

func NewRoom(roomname string) *Room {
	return &Room{roomname, 0, []*Client{}, &Broadcaster{}}
}
