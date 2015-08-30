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
    publish chan *Topic
    unpublish chan *Topic
}

var manager = Manager{
    topics: make(map[string] *Topic),
    subscribers: make(map[string] *Subscriber),
    publish: make(chan *Topic),
    unpublish: make(chan *Topic),
}


// Run manager
func RunManager() {
    defer manager.close()
    for {
        select {
        case topic := <- manager.publish:
            manager.topics[topic.id] = topic
            go topic.activate()
            log.Printf("new topic %v", topic.name)
            manager.broadcast(map[string] string { "tag": "publish", "id": topic.id, "name": topic.name})
        case topic := <- manager.unpublish:
            delete(manager.topics, topic.id)
            topic.close()
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
    manager.subscribers[sub.id] = sub
    log.Printf("New connection on %v", time.Now())
    go sub.writeMessage()
    sub.readMessage()
}

// broadcast data
func (m *Manager) broadcast(data custom) {
    result, _ := return_data(data)
    log.Printf("Broadcast data %v", data)
    for _, sub := range m.subscribers {
        sub.send <- result
    }
}

// Add new topic
func (m *Manager) addTopic(name string) *Topic{
    id := util.RandomID()
    topic := &Topic{
        id: id,
        name: name,
        subscribers: make(map[string] int),
        events: make(chan *Event),
        subscribe: make(chan *Subscriber),
        unsubscribe: make(chan *Subscriber),
    }
    m.publish <- topic
    return topic
}

// Close manager
func (m *Manager) close() {
    for _, topic := range m.topics {
        topic.close()
    }
    for _, sub := range m.subscribers {
        sub.close()
    }
    close(m.publish)
    close(m.unpublish)
}