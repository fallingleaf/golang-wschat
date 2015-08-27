package main

import (
    "flag"
    "net/http"
    "ocean/handlers"
    "ocean/pubsub"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
    flag.Parse()
    go pubsub.RunManager()
    http.Handle("/", handlers.InitRoutes())
    // serve static files
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
    http.ListenAndServe(*addr, nil)
}