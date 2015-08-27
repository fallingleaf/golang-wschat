package pubsub

import (
    "log"
    "time"
    "bytes"
    "encoding/json"
    "github.com/gorilla/websocket"
)


type Subscriber struct {
    id string
    ws * websocket.Conn
    send chan []byte
    // list of topics subscribed with status 
    topics map[string] int
}

const (
    // Time allowed to write a message to peer
    writeWait = 10 * time.Second
    
    // Time allowed to read the next pong message from peer
    pongWait = 60 * time.Second
    
    // Send ping to peer within this period. Must be less than pongWait
    pingPeriod = (pongWait * 9)/10
    
    // Maximum message size allowed
    maxMessageSize = 512
)

// Close websocket connection and set status
func (s *Subscriber) close() {
    s.ws.Close()
    event := &Event{
                    etype: e_status, 
                    source: s.id, 
                    data: status_offline,
                }
    for k := range s.topics {
        topic := manager.topics[k]
        topic.events <- event
    }
    manager.unregisters <- s
}

// Read websocket message
func (s *Subscriber) readMessage() {
    defer s.close()
    s.ws.SetReadLimit(maxMessageSize)
    s.ws.SetReadDeadline(time.Now().Add(pongWait))
    s.ws.SetPongHandler(func(string) error {
        s.ws.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })
    for {
        _, message, err := s.ws.ReadMessage()
        if err != nil {
            break
        }
        log.Printf("%v",string(message))
        // read Json data and push message to channel
        s.handleMessage(message)
    }
}

// Write data to websocket with the given message and type
func (s *Subscriber) write(mt int, payload []byte) error {
    s.ws.SetWriteDeadline(time.Now().Add(writeWait))
    return s.ws.WriteMessage(mt, payload)
}

// Write message to websocket
func (s *Subscriber) writeMessage() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        s.close()
    }()
    
    for {
        select {
        case m, ok := <- s.send:
            if !ok {
                s.write(websocket.CloseMessage, []byte{})
                return
            }
            if err := s.write(websocket.TextMessage, m); err != nil {
                return
            }
        case <-ticker.C:
            if err := s.write(websocket.PingMessage, []byte{}); err != nil {
                return
            }
        }
    }
}

// Handle incoming messages by type: chat, call
func (sub *Subscriber) handleMessage(message []byte) error {
    d := &Message{}
    decoder := json.NewDecoder(bytes.NewReader(message))
    err := decoder.Decode(d)
    if err != nil {
        log.Printf("Decoder error %v", err)
        return err
    }
    log.Printf("message %v", d)
    switch d.Type {
        case type_chat:
            topicId, _ := d.Data["topic"].(string)
            chatMessage, _ := d.Data["message"].(string)
            sub.sendChat(topicId, chatMessage)
        case type_method:
            sub.callMethod(d.Data)
        default:
            log.Printf("Unknown message type %v", d)
    }
    return nil
}

func (sub *Subscriber) sendChat(roomId, message string) {
    topic, ok := manager.topics[roomId]
    if !ok {
        log.Printf("Unknown topic %v", roomId)
        return
    }
    event := &Event{
                etype: e_message, 
                source: sub.id, 
                data: message,
            }
            
    topic.events <- event
}

func (sub *Subscriber) callMethod(data map_data) error {
    name, _ := data["api"].(string)
    handler, ok := methods[name]
    if !ok {
        log.Printf("Invalid method call: %v", name)
        return nil
    }
    return handler(data)
}
