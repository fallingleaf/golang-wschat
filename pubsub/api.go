package pubsub

import (
    "log"
    "net/http"
    "encoding/json"
)

type Result struct {
    Type    int `json:"type"`
    Code    int `json:"code,omitempty"`
    Success bool `json:"success,omitempty"`
    Data    custom `json:"data"`
    Message string `json:"message,omitempty"`
}

func return_data(data custom) ([]byte, error) {
    result := &Result {
        Type: type_push,
        Data: data,
    }
    
    return json.Marshal(result)
}

func return_success(data custom) ([]byte, error) {
    result := &Result {
        Type: type_method,
        Code: http.StatusOK,
        Success: true,
        Data: data,
        Message: "",
    }
    
    return json.Marshal(result)
}

func return_error(message string, code int) ([]byte, error) {
    result := &Result {
        Type: type_method,
        Code: code,
        Success: false,
        Data: "",
        Message: message,
    }
    return json.Marshal(result)
}

func return_message(topic *Topic, event *Event) ([]byte, error) {
    result := &Result {
        Type: type_chat,
        Data: map[string] string {
            "message": event.data.(string),
            "sender": event.source,
            "topic": topic.id,
        },
    }
    return json.Marshal(result)
}

func auth(sub *Subscriber, data map_data) []byte {
    log.Println("calling auth")
    return nil
}

func publish(sub *Subscriber, data map_data) []byte {
    name := data["topic"].(string)
    manager.addTopic(name)
    // auto register publisher
    //log.Println("topic added")
    //topic.subscribe <- sub
    result, _ := return_success("")
    return result
}

func subscribe(sub *Subscriber, data map_data) []byte {
    id := data["topic"].(string)
    topic, ok := manager.topics[id]
    if !ok {
        result, _ := return_error("Topic not found", http.StatusNotFound)
        return result
    }
    topic.subscribe <- sub
    result, _ := return_success("")
    return result
}

func syncTopic(sub *Subscriber, data map_data) []byte {
    topics := manager.topics
    var r []map[string]string
    for _, v := range topics {
        _, ok := sub.topics[v.id]
        if ok {
            continue
        }
        t := map[string] string { "id": v.id, "name": v.name}
        r = append(r, t)
    }
    j := map[string] custom { "tag": "sync", "topics": r}
    result, _ := return_data(j)
    return result
}

