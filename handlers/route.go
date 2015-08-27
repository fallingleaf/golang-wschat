package handlers

import (
    _ "net/http"
    "github.com/gorilla/mux"
)

// init all routes used in api
// return route instance
func InitRoutes() *mux.Router{
    r := mux.NewRouter()
    // add routes here
    r.HandleFunc("/", HomeHandler)
    r.HandleFunc("/ws", WebSocketHandler).Methods("GET")
    r.HandleFunc("/api/users", UserListHandler).Methods("GET", "POST")
    r.HandleFunc("/api/users/{id:[0-9]+}", UserDetailHandler).Methods("GET", "POST", "PUT")
    return r
}