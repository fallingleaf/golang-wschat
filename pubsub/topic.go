package pubsub

import (
    "log"
)

type Event struct {
    etype int
    source string
    data custom
}

type Topic struct {
    id string
    name string
    // list of subscribers with status
    subscribers map[string] int
    events chan *Event
    subscribe chan *Subscriber
    unsubscribe chan *Subscriber
}

func (topic *Topic) String() string {
    var t string = ""
    for k, _ := range topic.subscribers {
        t += "," + k
    }
    return t
}

func (topic *Topic) handleEvent(e *Event) {
    switch e.etype {
        case e_status:
            _, ok := topic.subscribers[e.source]
            if ok {
                topic.subscribers[e.source] = e.data.(int)
            }
        case e_message:
            result, _ := return_message(topic, e)
            for k, _ := range topic.subscribers {
                sub := manager.subscribers[k]
                sub.send <- result
            }
            
        default:
            log.Printf("Unknown event %v", e)
    }
}

func (topic *Topic) activate() {
    defer topic.close()
    for {
        select {
        case e := <- topic.events:
            topic.handleEvent(e)
        case sub := <- topic.subscribe:
            log.Printf("add new subsciber %v", sub.id)
            topic.subscribers[sub.id] = status_online
            sub.topics[topic.id] = status_online
        case sub := <- topic.unsubscribe:
            delete(topic.subscribers, sub.id)
        }
    }
    
}

func (topic *Topic) close() {
    for k, _ := range topic.subscribers {
        sub := manager.subscribers[k]
        delete(sub.topics, topic.id)
    }
    delete(manager.topics, topic.id)
    close(topic.events)
    close(topic.subscribe)
    close(topic.unsubscribe)
}