/*
 * Define user and chat models structure used for parsing 
 * real data from queries
 */
package database

type User struct {
    Id          int     `json:"id"`
    Username    string  `json:"username"`
    Password    string  `json:"-"`
    Token       string  `json:"token"`
}

type Room struct {
    id          int
    name        int
    creator_id  int
}

type Message struct {
    id          int
    message     string
    sender_id   int
    room_id     int
    timestamp   int
}