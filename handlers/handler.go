package handlers

import (
    "log"
    "net/http"
    "encoding/json"
    "github.com/gorilla/websocket"
    "ocean/pubsub"
    "text/template"
)

type custom interface{}

func writeResponse(w http.ResponseWriter, data custom) error {
    header := w.Header()
    header.Set("Content-Type", "application/json;charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    err := json.NewEncoder(w).Encode(data)
    return err
}

func sendResult(w http.ResponseWriter, data custom) error {
    result := map[string]interface{} {
        "success": true,
        "data": data,
    }
    return writeResponse(w, result)
}

func sendError(w http.ResponseWriter, message string, code int) error {
    result := map[string]custom {
        "success": false,
        "message": message,
        "code": code,
    }
    return writeResponse(w, result)
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
    var upgrader = websocket.Upgrader{
    	ReadBufferSize:  1024,
    	WriteBufferSize: 1024,
        CheckOrigin: func(req *http.Request) bool {return true},
    }
    
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
    
    pubsub.AddSubscriber(ws)
    
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
    homeTempl := template.Must(template.ParseFiles("templates/index.html"))
	homeTempl.Execute(w, r.Host)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    sendResult(w, "Hello")
}

func UserListHandler(w http.ResponseWriter, r *http.Request) {
    
}

func UserDetailHandler(w http.ResponseWriter, r *http.Request) {
}