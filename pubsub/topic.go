package pubsub

import (
    "log"
)

type Topic struct {
    id string
    name string
    // list of subscribers with status
    subscribers map[string] int
    events chan *Event
}

func (topic *Topic) handleEvent(e *Event) {
    switch e.etype {
        case e_status:
            _, ok := topic.subscribers[e.source]
            if ok {
                topic.subscribers[e.source] = e.data.(int)
            }
        case e_message:
            for k, _ := range topic.subscribers {
                sub := manager.subscribers[k]
                sub.send <- e.data.([]byte)
            }
        default:
            log.Printf("Unknown event %v", e)
    }
}

func (topic *Topic) activate() {
    for {
        select {
        case e := <- topic.events:
            topic.handleEvent(e)
        }
    }
}