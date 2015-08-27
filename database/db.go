package database

import (
    "log"
    "database/sql"
    _ "github.com/lib/pq"
)

var Dbconn *sql.DB

func Connect() {
    if Dbconn != nil {
        return
    }
    var err error
    Dbconn, err = sql.Open("postgres", "postgres://rosary:05021988@localhost/orion?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    //defer Dbconn.Close()
}