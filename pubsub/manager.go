package pubsub

import (
    "log"
    "time"
    "ocean/util"
    "github.com/gorilla/websocket"
)

type Manager struct {
    topics map[string] *Topic
    subscribers map[string] *Subscriber
    registers chan *Subscriber
    unregisters chan *Subscriber
}

var manager = Manager{
    topics: make(map[string] *Topic),
    subscribers: make(map[string] *Subscriber),
    registers: make(chan *Subscriber),
    unregisters: make(chan *Subscriber),
}


// Run manager
func RunManager() {
    for {
        select {
        case s := <- manager.registers:
            manager.subscribers[s.id] = s
        case s := <- manager.unregisters:
            delete(manager.subscribers, s.id)
            close(s.send)
        }
    }
}
 
// Add new subscriber with given id and websocket
func AddSubscriber(ws *websocket.Conn) {
    sub := &Subscriber{
            id: util.RandomID(), 
            send: make(chan []byte), 
            ws: ws, 
            topics: make(map[string] int),
        }
    manager.registers <- sub
    log.Printf("New connection on %v", time.Now())
    go sub.writeMessage()
    sub.readMessage()
}